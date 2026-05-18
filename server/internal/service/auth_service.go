package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/freelanceplatform/server/internal/config"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg           *config.Config
	userRepo      *repository.UserRepository
	tokenRepo     *repository.RefreshTokenRepository
	developerRepo *repository.DeveloperRepository
	clientRepo    *repository.ClientRepository

	resetCodes   map[string]string // key: emailOrPhone, value: 6-digit code
	resetCodesMu sync.RWMutex
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository, tokenRepo *repository.RefreshTokenRepository, developerRepo *repository.DeveloperRepository, clientRepo *repository.ClientRepository) *AuthService {
	return &AuthService{
		cfg:           cfg,
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		developerRepo: developerRepo,
		clientRepo:    clientRepo,
		resetCodes:    make(map[string]string),
	}
}

// Register creates a new user and returns the user along with an access/refresh token pair.
// Email and phone are mutually exclusive: exactly one must be provided.
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, string, string, error) {
	hasEmail := req.Email != nil && *req.Email != ""
	hasPhone := req.Phone != nil && *req.Phone != ""

	if !hasEmail && !hasPhone {
		return nil, "", "", errors.New("either email or phone is required")
	}
	if hasEmail && hasPhone {
		return nil, "", "", errors.New("provide either email or phone, not both")
	}

	// Check uniqueness
	if hasEmail {
		existing, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, "", "", fmt.Errorf("check email existence: %w", err)
		}
		if existing != nil {
			return nil, "", "", errors.New("email already registered")
		}
	}
	if hasPhone {
		existing, err := s.userRepo.GetByPhone(ctx, *req.Phone)
		if err != nil {
			return nil, "", "", fmt.Errorf("check phone existence: %w", err)
		}
		if existing != nil {
			return nil, "", "", errors.New("phone already registered")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", fmt.Errorf("hash password: %w", err)
	}

	emailVal := ""
	if hasEmail {
		emailVal = *req.Email
	}
	phoneVal := ""
	if hasPhone {
		phoneVal = *req.Phone
	}

	user, err := s.userRepo.Create(ctx, emailVal, phoneVal, string(hash), req.Nickname, req.UserType)
	if err != nil {
		return nil, "", "", fmt.Errorf("create user: %w", err)
	}

	// Auto-create developer or client profile so they appear in search immediately
	if user.UserType == "developer" || user.UserType == "both" {
		_, _ = s.developerRepo.Create(ctx, user.ID, map[string]interface{}{
			"availability": "available",
		})
	}
	if user.UserType == "client" || user.UserType == "both" {
		_, _ = s.clientRepo.Create(ctx, user.ID, map[string]interface{}{})
	}

	accessToken, refreshToken, err := s.GenerateTokenPair(user.ID, user.UserType)
	if err != nil {
		return nil, "", "", fmt.Errorf("generate token pair: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// Login authenticates a user and returns an access/refresh token pair.
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (string, string, error) {
	var user *model.User
	var err error

	hasEmail := req.Email != nil && *req.Email != ""
	hasPhone := req.Phone != nil && *req.Phone != ""

	if !hasEmail && !hasPhone {
		return "", "", errors.New("either email or phone is required")
	}
	if hasEmail && hasPhone {
		return "", "", errors.New("provide either email or phone, not both")
	}

	if hasEmail {
		user, err = s.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return "", "", fmt.Errorf("find user by email: %w", err)
		}
	} else {
		user, err = s.userRepo.GetByPhone(ctx, *req.Phone)
		if err != nil {
			return "", "", fmt.Errorf("find user by phone: %w", err)
		}
	}

	if user == nil {
		return "", "", errors.New("invalid credentials")
	}
	if user.Status == "deleted" {
		return "", "", errors.New("account has been deleted")
	}
	if user.Status == "suspended" {
		return "", "", errors.New("account has been suspended")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return "", "", fmt.Errorf("update last login: %w", err)
	}

	accessToken, refreshToken, err := s.GenerateTokenPair(user.ID, user.UserType)
	if err != nil {
		return "", "", fmt.Errorf("generate token pair: %w", err)
	}

	return accessToken, refreshToken, nil
}

// Logout revokes the given refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}
	if err := s.tokenRepo.Revoke(ctx, refreshToken); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

// RefreshToken validates the old refresh token, revokes it, and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	rt, err := s.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("get refresh token: %w", err)
	}
	if rt == nil {
		return "", "", errors.New("refresh token not found")
	}
	if rt.Revoked {
		return "", "", errors.New("refresh token has been revoked")
	}
	if time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("refresh token has expired")
	}

	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}
	if user.Status == "deleted" || user.Status == "suspended" {
		return "", "", errors.New("account is not active")
	}

	// Revoke old refresh token
	if err := s.tokenRepo.Revoke(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("revoke old refresh token: %w", err)
	}

	// Generate new token pair
	accessToken, newRefreshToken, err := s.GenerateTokenPair(user.ID, user.UserType)
	if err != nil {
		return "", "", fmt.Errorf("generate token pair: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// ChangePassword verifies the old password and updates to the new one.
func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	if oldPassword == "" || newPassword == "" {
		return errors.New("old password and new password are required")
	}
	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, string(hash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Revoke all refresh tokens to force re-login
	_ = s.tokenRepo.RevokeAllByUser(ctx, userID)

	return nil
}

// ForgotPassword generates a 6-digit verification code and stores it (MVP: in-memory map).
func (s *AuthService) ForgotPassword(ctx context.Context, emailOrPhone string) error {
	if emailOrPhone == "" {
		return errors.New("email or phone is required")
	}

	var user *model.User
	var err error
	if strings.Contains(emailOrPhone, "@") {
		user, err = s.userRepo.GetByEmail(ctx, emailOrPhone)
	} else {
		user, err = s.userRepo.GetByPhone(ctx, emailOrPhone)
	}
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		// Do not reveal whether the account exists for security
		return nil
	}

	code := generateResetCode()
	s.resetCodesMu.Lock()
	s.resetCodes[emailOrPhone] = code
	s.resetCodesMu.Unlock()

	// In production, send code via email/SMS. For MVP, it is logged/returned.
	_ = code // caller should retrieve via a separate endpoint or log
	return nil
}

// ResetPassword verifies the code and updates the password.
func (s *AuthService) ResetPassword(ctx context.Context, emailOrPhone, code, newPassword string) error {
	if emailOrPhone == "" || code == "" || newPassword == "" {
		return errors.New("email/phone, code, and new password are required")
	}
	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	s.resetCodesMu.RLock()
	storedCode, ok := s.resetCodes[emailOrPhone]
	s.resetCodesMu.RUnlock()

	if !ok || storedCode != code {
		return errors.New("invalid or expired verification code")
	}

	var user *model.User
	var err error
	if strings.Contains(emailOrPhone, "@") {
		user, err = s.userRepo.GetByEmail(ctx, emailOrPhone)
	} else {
		user, err = s.userRepo.GetByPhone(ctx, emailOrPhone)
	}
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Clean up the used code
	s.resetCodesMu.Lock()
	delete(s.resetCodes, emailOrPhone)
	s.resetCodesMu.Unlock()

	// Revoke all refresh tokens
	_ = s.tokenRepo.RevokeAllByUser(ctx, user.ID)

	return nil
}

// ValidateToken parses and validates a JWT access token, returning its claims.
func (s *AuthService) ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateTokenPair creates an access token (15 min) and stores a refresh token (7 days) in DB.
func (s *AuthService) GenerateTokenPair(userID, userType string) (string, string, error) {
	accessExpiry := time.Duration(s.cfg.JWT.AccessExpiry) * time.Minute
	if accessExpiry == 0 {
		accessExpiry = 15 * time.Minute
	}

	now := time.Now()
	accessClaims := model.JWTClaims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.JWT.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiry)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	// Refresh token: random UUID stored in DB
	refreshTokenString := uuid.New().String()
	refreshExpiryDays := s.cfg.JWT.RefreshExpiry
	if refreshExpiryDays == 0 {
		refreshExpiryDays = 7
	}
	expiresAt := now.AddDate(0, 0, refreshExpiryDays)

	ctx := context.Background()
	_, err = s.tokenRepo.Create(ctx, userID, refreshTokenString, nil, nil, expiresAt)
	if err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// VerifyEmail marks the user's email as verified.
func (s *AuthService) VerifyEmail(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.EmailVerified {
		return errors.New("email already verified")
	}
	return s.userRepo.VerifyEmail(ctx, userID)
}

// VerifyPhone marks the user's phone as verified.
func (s *AuthService) VerifyPhone(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.PhoneVerified {
		return errors.New("phone already verified")
	}
	return s.userRepo.VerifyPhone(ctx, userID)
}

// GetResetCode returns the stored reset code for a given email/phone (for MVP testing).
func (s *AuthService) GetResetCode(emailOrPhone string) string {
	s.resetCodesMu.RLock()
	defer s.resetCodesMu.RUnlock()
	return s.resetCodes[emailOrPhone]
}

// generateResetCode creates a random 6-digit code.
func generateResetCode() string {
	t := time.Now().UnixNano()
	return fmt.Sprintf("%06d", t%1000000)
}

// GetUserIDFromToken is a helper to extract userID from a token string.
func (s *AuthService) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ParseResetCodeToInt is a utility for callers that need the code as int.
func ParseResetCodeToInt(code string) (int, error) {
	return strconv.Atoi(code)
}

// IsResetCodeValid checks the code without consuming it.
func (s *AuthService) IsResetCodeValid(emailOrPhone, code string) bool {
	s.resetCodesMu.RLock()
	defer s.resetCodesMu.RUnlock()
	stored, ok := s.resetCodes[emailOrPhone]
	return ok && stored == code
}

// RevokeAllUserTokens revokes all refresh tokens for a user.
func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.tokenRepo.RevokeAllByUser(ctx, userID)
}

// Ensure unused import hint
var _ = strings.TrimSpace
var _ = strconv.Itoa
