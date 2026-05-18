package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type DeveloperRepository struct {
	db *DB
}

func NewDeveloperRepository(db *DB) *DeveloperRepository {
	return &DeveloperRepository{db: db}
}

// scanDeveloperProfile scans a single developer profile row into a model.DeveloperProfile.
func scanDeveloperProfile(row pgx.Row) (*model.DeveloperProfile, error) {
	var dp model.DeveloperProfile
	err := row.Scan(
		&dp.ID, &dp.UserID, &dp.RealName, &dp.Title, &dp.Bio,
		&dp.HourlyRate, &dp.Availability, &dp.ExperienceYears, &dp.Location,
		&dp.GithubURL, &dp.LinkedinURL, &dp.WebsiteURL,
		&dp.Verified, &dp.VerifiedAt,
		&dp.RatingAvg, &dp.RatingCount, &dp.TotalEarnings, &dp.CompletedProjects,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *DeveloperRepository) Create(ctx context.Context, userID string, req map[string]interface{}) (*model.DeveloperProfile, error) {
	var dp model.DeveloperProfile
	query := `
		INSERT INTO developer_profiles (user_id, real_name, title, bio, hourly_rate,
			availability, experience_years, location, github_url, linkedin_url, website_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, real_name, title, bio, hourly_rate, availability,
			experience_years, location, github_url, linkedin_url, website_url,
			verified, verified_at, rating_avg, rating_count, total_earnings,
			completed_projects, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		userID,
		fieldPtr[string](req, "real_name"),
		fieldPtr[string](req, "title"),
		fieldPtr[string](req, "bio"),
		fieldPtr[float64](req, "hourly_rate"),
		fieldStr(req, "availability", "available"),
		fieldInt(req, "experience_years", 0),
		fieldPtr[string](req, "location"),
		fieldPtr[string](req, "github_url"),
		fieldPtr[string](req, "linkedin_url"),
		fieldPtr[string](req, "website_url"),
	).Scan(
		&dp.ID, &dp.UserID, &dp.RealName, &dp.Title, &dp.Bio,
		&dp.HourlyRate, &dp.Availability, &dp.ExperienceYears, &dp.Location,
		&dp.GithubURL, &dp.LinkedinURL, &dp.WebsiteURL,
		&dp.Verified, &dp.VerifiedAt,
		&dp.RatingAvg, &dp.RatingCount, &dp.TotalEarnings, &dp.CompletedProjects,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create developer profile: %w", err)
	}
	return &dp, nil
}

func (r *DeveloperRepository) GetByID(ctx context.Context, id string) (*model.DeveloperProfile, error) {
	query := `
		SELECT id, user_id, real_name, title, bio, hourly_rate, availability,
			experience_years, location, github_url, linkedin_url, website_url,
			verified, verified_at, rating_avg, rating_count, total_earnings,
			completed_projects, created_at, updated_at
		FROM developer_profiles
		WHERE id = $1
	`
	dp, err := scanDeveloperProfile(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get developer profile by id: %w", err)
	}
	return dp, nil
}

func (r *DeveloperRepository) GetByUserID(ctx context.Context, userID string) (*model.DeveloperProfile, error) {
	query := `
		SELECT id, user_id, real_name, title, bio, hourly_rate, availability,
			experience_years, location, github_url, linkedin_url, website_url,
			verified, verified_at, rating_avg, rating_count, total_earnings,
			completed_projects, created_at, updated_at
		FROM developer_profiles
		WHERE user_id = $1
	`
	dp, err := scanDeveloperProfile(r.db.Pool.QueryRow(ctx, query, userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get developer profile by user id: %w", err)
	}
	return dp, nil
}

func (r *DeveloperRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*model.DeveloperProfile, error) {
	if len(fields) == 0 {
		return r.GetByID(ctx, id)
	}

	allowed := map[string]bool{
		"real_name": true, "title": true, "bio": true, "hourly_rate": true,
		"availability": true, "experience_years": true, "location": true,
		"github_url": true, "linkedin_url": true, "website_url": true,
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
		UPDATE developer_profiles SET %s
		WHERE id = $%d
		RETURNING id, user_id, real_name, title, bio, hourly_rate, availability,
			experience_years, location, github_url, linkedin_url, website_url,
			verified, verified_at, rating_avg, rating_count, total_earnings,
			completed_projects, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIdx)

	dp, err := scanDeveloperProfile(r.db.Pool.QueryRow(ctx, query, args...))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("developer profile not found")
		}
		return nil, fmt.Errorf("update developer profile: %w", err)
	}
	return dp, nil
}

// ─── Skills ───

func (r *DeveloperRepository) AddSkill(ctx context.Context, developerID, skillName, proficiency string, yearsExp *int) (*model.DeveloperSkill, error) {
	var s model.DeveloperSkill
	query := `
		INSERT INTO developer_skills (developer_id, skill_name, proficiency, years_experience)
		VALUES ($1, $2, $3, $4)
		RETURNING id, developer_id, skill_name, proficiency, years_experience, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		developerID, skillName, proficiency, yearsExp,
	).Scan(
		&s.ID, &s.DeveloperID, &s.SkillName, &s.Proficiency, &s.YearsExperience, &s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("add developer skill: %w", err)
	}
	return &s, nil
}

func (r *DeveloperRepository) UpdateSkill(ctx context.Context, id string, proficiency string, yearsExp *int) error {
	query := `UPDATE developer_skills SET proficiency = $2, years_experience = $3 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, proficiency, yearsExp)
	if err != nil {
		return fmt.Errorf("update developer skill: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("developer skill not found")
	}
	return nil
}

func (r *DeveloperRepository) DeleteSkill(ctx context.Context, id string) error {
	query := `DELETE FROM developer_skills WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete developer skill: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("developer skill not found")
	}
	return nil
}

func (r *DeveloperRepository) ListSkills(ctx context.Context, developerID string) ([]*model.DeveloperSkill, error) {
	query := `
		SELECT id, developer_id, skill_name, proficiency, years_experience, created_at
		FROM developer_skills
		WHERE developer_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.Pool.Query(ctx, query, developerID)
	if err != nil {
		return nil, fmt.Errorf("list developer skills: %w", err)
	}
	defer rows.Close()

	var skills []*model.DeveloperSkill
	for rows.Next() {
		var s model.DeveloperSkill
		if err := rows.Scan(
			&s.ID, &s.DeveloperID, &s.SkillName, &s.Proficiency, &s.YearsExperience, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan developer skill: %w", err)
		}
		skills = append(skills, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate developer skills: %w", err)
	}
	return skills, nil
}

// ─── Portfolio ───

func (r *DeveloperRepository) AddPortfolio(ctx context.Context, developerID string, p *model.DeveloperPortfolio) (*model.DeveloperPortfolio, error) {
	imageURLs, _ := json.Marshal(p.ImageURLs)
	if imageURLs == nil {
		imageURLs = []byte("[]")
	}
	techStack, _ := json.Marshal(p.TechStack)
	if techStack == nil {
		techStack = []byte("[]")
	}

	var dp model.DeveloperPortfolio
	query := `
		INSERT INTO developer_portfolio (developer_id, title, description, project_url,
			image_urls, tech_stack, start_date, end_date, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, developer_id, title, description, project_url,
			image_urls, tech_stack, start_date, end_date, sort_order, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		developerID, p.Title, p.Description, p.ProjectURL,
		imageURLs, techStack, p.StartDate, p.EndDate, p.SortOrder,
	).Scan(
		&dp.ID, &dp.DeveloperID, &dp.Title, &dp.Description, &dp.ProjectURL,
		&dp.ImageURLs, &dp.TechStack, &dp.StartDate, &dp.EndDate, &dp.SortOrder,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("add developer portfolio: %w", err)
	}
	return &dp, nil
}

func (r *DeveloperRepository) UpdatePortfolio(ctx context.Context, id string, p *model.DeveloperPortfolio) (*model.DeveloperPortfolio, error) {
	imageURLs, _ := json.Marshal(p.ImageURLs)
	if imageURLs == nil {
		imageURLs = []byte("[]")
	}
	techStack, _ := json.Marshal(p.TechStack)
	if techStack == nil {
		techStack = []byte("[]")
	}

	var dp model.DeveloperPortfolio
	query := `
		UPDATE developer_portfolio SET
			title = $2, description = $3, project_url = $4,
			image_urls = $5, tech_stack = $6, start_date = $7,
			end_date = $8, sort_order = $9
		WHERE id = $1
		RETURNING id, developer_id, title, description, project_url,
			image_urls, tech_stack, start_date, end_date, sort_order, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		id, p.Title, p.Description, p.ProjectURL,
		imageURLs, techStack, p.StartDate, p.EndDate, p.SortOrder,
	).Scan(
		&dp.ID, &dp.DeveloperID, &dp.Title, &dp.Description, &dp.ProjectURL,
		&dp.ImageURLs, &dp.TechStack, &dp.StartDate, &dp.EndDate, &dp.SortOrder,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("portfolio not found")
		}
		return nil, fmt.Errorf("update developer portfolio: %w", err)
	}
	return &dp, nil
}

func (r *DeveloperRepository) DeletePortfolio(ctx context.Context, id string) error {
	query := `DELETE FROM developer_portfolio WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete developer portfolio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("portfolio not found")
	}
	return nil
}

func (r *DeveloperRepository) ListPortfolio(ctx context.Context, developerID string) ([]*model.DeveloperPortfolio, error) {
	query := `
		SELECT id, developer_id, title, description, project_url,
			image_urls, tech_stack, start_date, end_date, sort_order, created_at, updated_at
		FROM developer_portfolio
		WHERE developer_id = $1
		ORDER BY sort_order, created_at
	`
	rows, err := r.db.Pool.Query(ctx, query, developerID)
	if err != nil {
		return nil, fmt.Errorf("list developer portfolio: %w", err)
	}
	defer rows.Close()

	var items []*model.DeveloperPortfolio
	for rows.Next() {
		var dp model.DeveloperPortfolio
		if err := rows.Scan(
			&dp.ID, &dp.DeveloperID, &dp.Title, &dp.Description, &dp.ProjectURL,
			&dp.ImageURLs, &dp.TechStack, &dp.StartDate, &dp.EndDate, &dp.SortOrder,
			&dp.CreatedAt, &dp.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan developer portfolio: %w", err)
		}
		items = append(items, &dp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate developer portfolio: %w", err)
	}
	return items, nil
}

// ─── Search ───

func (r *DeveloperRepository) Search(ctx context.Context, skill string, minRate, maxRate *float64, availability string, page, pageSize int) ([]*model.DeveloperProfile, int64, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if skill != "" {
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM developer_skills WHERE developer_skills.developer_id = developer_profiles.id AND skill_name ILIKE $%d)`, argIdx,
		))
		args = append(args, "%"+skill+"%")
		argIdx++
	}

	if minRate != nil {
		conditions = append(conditions, fmt.Sprintf(`hourly_rate >= $%d`, argIdx))
		args = append(args, *minRate)
		argIdx++
	}

	if maxRate != nil {
		conditions = append(conditions, fmt.Sprintf(`hourly_rate <= $%d`, argIdx))
		args = append(args, *maxRate)
		argIdx++
	}

	if availability != "" {
		conditions = append(conditions, fmt.Sprintf(`availability = $%d`, argIdx))
		args = append(args, availability)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM developer_profiles %s`, whereClause)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count developer profiles: %w", err)
	}

	// Fetch page
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, user_id, real_name, title, bio, hourly_rate, availability,
			experience_years, location, github_url, linkedin_url, website_url,
			verified, verified_at, rating_avg, rating_count, total_earnings,
			completed_projects, created_at, updated_at
		FROM developer_profiles
		%s
		ORDER BY rating_avg DESC NULLS LAST, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, pageSize, offset)
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("search developer profiles: %w", err)
	}
	defer rows.Close()

	var results []*model.DeveloperProfile
	for rows.Next() {
		dp, scanErr := scanDeveloperProfileFromRows(rows)
		if scanErr != nil {
			return nil, 0, fmt.Errorf("scan developer profile: %w", scanErr)
		}
		results = append(results, dp)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate developer profiles: %w", err)
	}
	return results, total, nil
}

// ─── Stats ───

func (r *DeveloperRepository) UpdateRating(ctx context.Context, id string, avgRating float64, count int) error {
	query := `UPDATE developer_profiles SET rating_avg = $2, rating_count = $3 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, avgRating, count)
	if err != nil {
		return fmt.Errorf("update developer rating: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("developer profile not found")
	}
	return nil
}

func (r *DeveloperRepository) UpdateEarnings(ctx context.Context, id string, amount float64) error {
	query := `UPDATE developer_profiles SET total_earnings = total_earnings + $2 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, amount)
	if err != nil {
		return fmt.Errorf("update developer earnings: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("developer profile not found")
	}
	return nil
}

// ─── Helpers ───

// scanDeveloperProfileFromRows scans from pgx.Rows (used in multi-row queries).
func scanDeveloperProfileFromRows(rows pgx.Rows) (*model.DeveloperProfile, error) {
	var dp model.DeveloperProfile
	err := rows.Scan(
		&dp.ID, &dp.UserID, &dp.RealName, &dp.Title, &dp.Bio,
		&dp.HourlyRate, &dp.Availability, &dp.ExperienceYears, &dp.Location,
		&dp.GithubURL, &dp.LinkedinURL, &dp.WebsiteURL,
		&dp.Verified, &dp.VerifiedAt,
		&dp.RatingAvg, &dp.RatingCount, &dp.TotalEarnings, &dp.CompletedProjects,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

// fieldPtr extracts an optional field from a map as a typed pointer.
func fieldPtr[T any](m map[string]interface{}, key string) *T {
	val, ok := m[key]
	if !ok || val == nil {
		return nil
	}
	typed, ok := val.(T)
	if !ok {
		return nil
	}
	return &typed
}

// fieldStr extracts a string field from a map with a default value.
func fieldStr(m map[string]interface{}, key string, defaultVal string) string {
	val, ok := m[key]
	if !ok || val == nil {
		return defaultVal
	}
	s, ok := val.(string)
	if !ok {
		return defaultVal
	}
	return s
}

// fieldInt extracts an int field from a map with a default value.
func fieldInt(m map[string]interface{}, key string, defaultVal int) int {
	val, ok := m[key]
	if !ok || val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultVal
	}
}
