package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	db *DB
}

func NewClientRepository(db *DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// scanClientProfile scans a single client profile row into a model.ClientProfile.
func scanClientProfile(row pgx.Row) (*model.ClientProfile, error) {
	var cp model.ClientProfile
	err := row.Scan(
		&cp.ID, &cp.UserID, &cp.CompanyName, &cp.CompanyLogoURL, &cp.CompanyWebsite,
		&cp.Industry, &cp.CompanySize,
		&cp.Verified, &cp.VerifiedAt,
		&cp.TotalSpent, &cp.PostedProjects,
		&cp.CreatedAt, &cp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *ClientRepository) Create(ctx context.Context, userID string, req map[string]interface{}) (*model.ClientProfile, error) {
	var cp model.ClientProfile
	query := `
		INSERT INTO client_profiles (user_id, company_name, company_logo_url, company_website,
			industry, company_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, company_name, company_logo_url, company_website,
			industry, company_size, verified, verified_at, total_spent, posted_projects,
			created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		userID,
		fieldPtr[string](req, "company_name"),
		fieldPtr[string](req, "company_logo_url"),
		fieldPtr[string](req, "company_website"),
		fieldPtr[string](req, "industry"),
		fieldPtr[string](req, "company_size"),
	).Scan(
		&cp.ID, &cp.UserID, &cp.CompanyName, &cp.CompanyLogoURL, &cp.CompanyWebsite,
		&cp.Industry, &cp.CompanySize,
		&cp.Verified, &cp.VerifiedAt,
		&cp.TotalSpent, &cp.PostedProjects,
		&cp.CreatedAt, &cp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create client profile: %w", err)
	}
	return &cp, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*model.ClientProfile, error) {
	query := `
		SELECT id, user_id, company_name, company_logo_url, company_website,
			industry, company_size, verified, verified_at, total_spent, posted_projects,
			created_at, updated_at
		FROM client_profiles
		WHERE id = $1
	`
	cp, err := scanClientProfile(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get client profile by id: %w", err)
	}
	return cp, nil
}

func (r *ClientRepository) GetByUserID(ctx context.Context, userID string) (*model.ClientProfile, error) {
	query := `
		SELECT id, user_id, company_name, company_logo_url, company_website,
			industry, company_size, verified, verified_at, total_spent, posted_projects,
			created_at, updated_at
		FROM client_profiles
		WHERE user_id = $1
	`
	cp, err := scanClientProfile(r.db.Pool.QueryRow(ctx, query, userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get client profile by user id: %w", err)
	}
	return cp, nil
}

func (r *ClientRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*model.ClientProfile, error) {
	if len(fields) == 0 {
		return r.GetByID(ctx, id)
	}

	allowed := map[string]bool{
		"company_name": true, "company_logo_url": true, "company_website": true,
		"industry": true, "company_size": true,
	}

	var setClauses []string
	var args []interface{}
	argIdx := 1

	for col, val := range fields {
		if !allowed[col] {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE client_profiles SET %s
		WHERE id = $%d
		RETURNING id, user_id, company_name, company_logo_url, company_website,
			industry, company_size, verified, verified_at, total_spent, posted_projects,
			created_at, updated_at
	`, strings.Join(setClauses, ", "), argIdx)

	cp, err := scanClientProfile(r.db.Pool.QueryRow(ctx, query, args...))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("client profile not found")
		}
		return nil, fmt.Errorf("update client profile: %w", err)
	}
	return cp, nil
}

func (r *ClientRepository) Verify(ctx context.Context, id string) error {
	query := `UPDATE client_profiles SET verified = true, verified_at = now() WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("verify client profile: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("client profile not found")
	}
	return nil
}
