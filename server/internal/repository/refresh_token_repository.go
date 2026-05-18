package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepository struct {
	db *DB
}

func NewRefreshTokenRepository(db *DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID, token string, deviceInfo, ipAddress *string, expiresAt time.Time) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `
		INSERT INTO refresh_tokens (user_id, token, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, token, device_info, ip_address, expires_at, revoked, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		userID, token, deviceInfo, ipAddress, expiresAt,
	).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.DeviceInfo, &rt.IPAddress,
		&rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `
		SELECT id, user_id, token, device_info, ip_address, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.DeviceInfo, &rt.IPAddress,
		&rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get refresh token by token: %w", err)
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`
	tag, err := r.db.Pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("refresh token not found")
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1 AND revoked = false`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens for user: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < now() OR revoked = true`
	_, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("delete expired refresh tokens: %w", err)
	}
	return nil
}
