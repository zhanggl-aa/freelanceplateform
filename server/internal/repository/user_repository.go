package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, phone, passwordHash, nickname, userType string) (*model.User, error) {
	var u model.User
	query := `
		INSERT INTO users (email, phone, password_hash, nickname, user_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		nilIfEmpty(email), nilIfEmpty(phone), passwordHash, nickname, userType,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND status != 'deleted'
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND status != 'deleted'
	`
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE phone = $1 AND status != 'deleted'
	`
	err := r.db.Pool.QueryRow(ctx, query, phone).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by phone: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByWechatOpenID(ctx context.Context, openid string) (*model.User, error) {
	var u model.User
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE wechat_openid = $1 AND status != 'deleted'
	`
	err := r.db.Pool.QueryRow(ctx, query, openid).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by wechat openid: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	var u model.User
	query := `
		UPDATE users SET
			email = $2,
			phone = $3,
			avatar_url = $4,
			nickname = $5,
			user_type = $6,
			status = $7
		WHERE id = $1 AND status != 'deleted'
		RETURNING id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.Phone,
		user.AvatarURL,
		user.Nickname,
		user.UserType,
		user.Status,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
		&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found or deleted")
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	query := `UPDATE users SET last_login_at = now() WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or deleted")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2 WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or deleted")
	}
	return nil
}

func (r *UserRepository) VerifyEmail(ctx context.Context, id string) error {
	query := `UPDATE users SET email_verified = true WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("verify email: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or deleted")
	}
	return nil
}

func (r *UserRepository) VerifyPhone(ctx context.Context, id string) error {
	query := `UPDATE users SET phone_verified = true WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("verify phone: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or deleted")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET status = 'deleted' WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE status != 'deleted'`
	if err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.WechatOpenID, &u.WechatUnionID,
			&u.AvatarURL, &u.Nickname, &u.UserType, &u.Status, &u.EmailVerified, &u.PhoneVerified,
			&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users: %w", err)
	}
	return users, total, nil
}

// nilIfEmpty returns nil for empty strings so that nullable columns are set to NULL.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
