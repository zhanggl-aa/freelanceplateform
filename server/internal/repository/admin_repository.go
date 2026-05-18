package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type AdminRepository struct {
	db *DB
}

func NewAdminRepository(db *DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetByUserID(ctx context.Context, userID string) (*model.AdminRole, error) {
	var a model.AdminRole
	query := `
		SELECT id, user_id, role, permissions, created_at, updated_at
		FROM admin_roles
		WHERE user_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&a.ID, &a.UserID, &a.Role, &a.Permissions, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get admin role by user id: %w", err)
	}
	return &a, nil
}

func (r *AdminRepository) IsAdmin(ctx context.Context, userID string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM admin_roles WHERE user_id = $1)`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *AdminRepository) Create(ctx context.Context, userID, role string, permissions *string) (*model.AdminRole, error) {
	var a model.AdminRole
	query := `
		INSERT INTO admin_roles (user_id, role, permissions)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, role, permissions, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, userID, role, permissions).Scan(
		&a.ID, &a.UserID, &a.Role, &a.Permissions, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create admin role: %w", err)
	}
	return &a, nil
}

func (r *AdminRepository) DashboardStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var userCount, projectCount, activeContracts, pendingDisputes int64
	var totalRevenue float64

	queries := []struct {
		sql    string
		target *interface{}
	}{
		{`SELECT COUNT(*) FROM users WHERE status != 'deleted'`, nil},
		{`SELECT COUNT(*) FROM projects`, nil},
		{`SELECT COUNT(*) FROM contracts WHERE status = 'active'`, nil},
		{`SELECT COALESCE(SUM(platform_fee), 0) FROM payments WHERE status IN ('released')`, nil},
		{`SELECT COUNT(*) FROM disputes WHERE status IN ('open', 'under_review')`, nil},
	}

	type intPtr = *int64
	type floatPtr = *float64

	// Execute each query
	if err := r.db.Pool.QueryRow(ctx, queries[0].sql).Scan(&userCount); err != nil {
		return nil, fmt.Errorf("get user count: %w", err)
	}
	if err := r.db.Pool.QueryRow(ctx, queries[1].sql).Scan(&projectCount); err != nil {
		return nil, fmt.Errorf("get project count: %w", err)
	}
	if err := r.db.Pool.QueryRow(ctx, queries[2].sql).Scan(&activeContracts); err != nil {
		return nil, fmt.Errorf("get active contracts: %w", err)
	}
	if err := r.db.Pool.QueryRow(ctx, queries[3].sql).Scan(&totalRevenue); err != nil {
		return nil, fmt.Errorf("get total revenue: %w", err)
	}
	if err := r.db.Pool.QueryRow(ctx, queries[4].sql).Scan(&pendingDisputes); err != nil {
		return nil, fmt.Errorf("get pending disputes: %w", err)
	}

	stats["user_count"] = userCount
	stats["project_count"] = projectCount
	stats["active_contracts"] = activeContracts
	stats["total_revenue"] = totalRevenue
	stats["pending_disputes"] = pendingDisputes

	return stats, nil
}

func (r *AdminRepository) ListUsers(ctx context.Context, search, status string, page, pageSize int) ([]*model.User, int64, error) {
	var total int64
	var countArgs []interface{}
	countQuery := `SELECT COUNT(*) FROM users WHERE status != 'deleted'`
	argIdx := 1

	var conditions []string
	if search != "" {
		conditions = append(conditions, fmt.Sprintf(`(nickname ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)`, argIdx, argIdx, argIdx))
		countArgs = append(countArgs, "%"+search+"%")
		argIdx++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf(`status = $%d`, argIdx))
		countArgs = append(countArgs, status)
		argIdx++
	}

	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	if err := r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, email, phone, password_hash, wechat_openid, wechat_unionid,
			avatar_url, nickname, user_type, status, email_verified, phone_verified,
			last_login_at, created_at, updated_at
		FROM users
		WHERE status != 'deleted'
	`
	var queryArgs []interface{}
	queryArgIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND (nickname ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)`, queryArgIdx, queryArgIdx, queryArgIdx)
		queryArgs = append(queryArgs, "%"+search+"%")
		queryArgIdx++
	}
	if status != "" {
		query += fmt.Sprintf(` AND status = $%d`, queryArgIdx)
		queryArgs = append(queryArgs, status)
		queryArgIdx++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, queryArgIdx, queryArgIdx+1)
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, query, queryArgs...)
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

func (r *AdminRepository) UpdateUserStatus(ctx context.Context, userID, status string) error {
	query := `UPDATE users SET status = $2 WHERE id = $1 AND status != 'deleted'`
	tag, err := r.db.Pool.Exec(ctx, query, userID, status)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or deleted")
	}
	return nil
}
