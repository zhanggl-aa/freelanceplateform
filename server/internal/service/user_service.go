package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type UserService struct {
	userRepo      *repository.UserRepository
	developerRepo *repository.DeveloperRepository
	clientRepo    *repository.ClientRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	developerRepo *repository.DeveloperRepository,
	clientRepo *repository.ClientRepository,
) *UserService {
	return &UserService{
		userRepo:      userRepo,
		developerRepo: developerRepo,
		clientRepo:    clientRepo,
	}
}

// GetCurrentUser returns the user along with any developer and client profiles.
func (s *UserService) GetCurrentUser(ctx context.Context, userID string) (*model.User, *model.DeveloperProfile, *model.ClientProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user: %w", err)
	}

	var devProfile *model.DeveloperProfile
	var clientProfile *model.ClientProfile

	// Attempt to load developer profile if user type includes developer
	if user.UserType == "developer" || user.UserType == "both" {
		devProfile, err = s.developerRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get developer profile: %w", err)
		}
	}

	// Attempt to load client profile if user type includes client
	if user.UserType == "client" || user.UserType == "both" {
		clientProfile, err = s.clientRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get client profile: %w", err)
		}
	}

	return user, devProfile, clientProfile, nil
}

// GetUserByID returns a user's public profile.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	// Clear sensitive fields for public profile
	user.PasswordHash = ""
	user.WechatOpenID = nil
	user.WechatUnionID = nil

	return user, nil
}

// UpdateProfile updates a user's nickname and avatar.
func (s *UserService) UpdateProfile(ctx context.Context, userID, nickname, avatarURL string) (*model.User, error) {
	if nickname == "" {
		return nil, errors.New("nickname is required")
	}
	if len(nickname) < 2 || len(nickname) > 50 {
		return nil, errors.New("nickname must be between 2 and 50 characters")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Nickname = nickname
	if avatarURL != "" {
		user.AvatarURL = &avatarURL
	}

	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return updated, nil
}

// DeleteAccount performs a soft delete of the user account.
func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.Status == "deleted" {
		return errors.New("account already deleted")
	}

	return s.userRepo.Delete(ctx, userID)
}
