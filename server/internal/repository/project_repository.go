package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type ProjectRepository struct {
	db *DB
}

func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

const projectColumns = `
	p.id::text, p.client_id::text, p.category_id::text, p.title, p.description,
	p.budget_min, p.budget_max, p.budget_type, p.deadline, p.tech_stack,
	p.status, p.view_count, p.bookmark_count, p.bid_count, p.bid_deadline,
	p.assigned_developer_id::text, p.featured, p.created_at, p.updated_at,
	u.nickname, u.avatar_url, c.name
`

const projectJoin = `
	FROM projects p
	JOIN users u ON u.id = p.client_id
	JOIN project_categories c ON c.id = p.category_id
`

func scanProject(row pgx.Row) (*model.Project, error) {
	p := &model.Project{}
	var clientName *string
	var clientAvatar *string
	var categoryName *string

	err := row.Scan(
		&p.ID, &p.ClientID, &p.CategoryID, &p.Title, &p.Description,
		&p.BudgetMin, &p.BudgetMax, &p.BudgetType, &p.Deadline, &p.TechStack,
		&p.Status, &p.ViewCount, &p.BookmarkCount, &p.BidCount, &p.BidDeadline,
		&p.AssignedDeveloperID, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
		&clientName, &clientAvatar, &categoryName,
	)
	if err != nil {
		return nil, err
	}
	p.ClientName = clientName
	p.ClientAvatar = clientAvatar
	p.CategoryName = categoryName
	return p, nil
}

func (r *ProjectRepository) Create(ctx context.Context, clientID, categoryID, title, description string, budgetMin, budgetMax *float64, budgetType string, deadline *string, techStack []string) (*model.Project, error) {
	query := `
		INSERT INTO projects (client_id, category_id, title, description, budget_min, budget_max, budget_type, deadline, tech_stack)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text, client_id::text, category_id::text, title, description,
			budget_min, budget_max, budget_type, deadline, tech_stack,
			status, view_count, bookmark_count, bid_count, bid_deadline,
			assigned_developer_id::text, featured, created_at, updated_at
	`
	p := &model.Project{}
	err := r.db.Pool.QueryRow(ctx, query,
		clientID, categoryID, title, description, budgetMin, budgetMax, budgetType, deadline, techStack,
	).Scan(
		&p.ID, &p.ClientID, &p.CategoryID, &p.Title, &p.Description,
		&p.BudgetMin, &p.BudgetMax, &p.BudgetType, &p.Deadline, &p.TechStack,
		&p.Status, &p.ViewCount, &p.BookmarkCount, &p.BidCount, &p.BidDeadline,
		&p.AssignedDeveloperID, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	query := `SELECT ` + projectColumns + projectJoin + ` WHERE p.id = $1`
	p, err := scanProject(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"title": true, "description": true, "budget_min": true, "budget_max": true,
		"budget_type": true, "deadline": true, "tech_stack": true, "status": true,
		"bid_deadline": true, "assigned_developer_id": true, "featured": true,
		"category_id": true,
	}

	setClauses := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+1)
	argIdx := 1

	for col, val := range fields {
		if !allowed[col] {
			return fmt.Errorf("invalid column: %s", col)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE projects SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found: %s", id)
	}
	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1 AND status = 'draft'`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found or not in draft status: %s", id)
	}
	return nil
}

func (r *ProjectRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE projects SET status = $1 WHERE id = $2`
	tag, err := r.db.Pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update project status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found: %s", id)
	}
	return nil
}

func (r *ProjectRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `UPDATE projects SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("increment view count: %w", err)
	}
	return nil
}

func (r *ProjectRepository) IncrementBidCount(ctx context.Context, id string) error {
	query := `UPDATE projects SET bid_count = bid_count + 1 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("increment bid count: %w", err)
	}
	return nil
}

func (r *ProjectRepository) Search(ctx context.Context, categoryID, status, keyword *string, minBudget, maxBudget *float64, techStack []string, page, pageSize int) ([]*model.Project, int64, error) {
	conditions := []string{"1=1"}
	args := make([]interface{}, 0)
	argIdx := 1

	if categoryID != nil && *categoryID != "" {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d::uuid", argIdx))
		args = append(args, *categoryID)
		argIdx++
	}
	if status != nil && *status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, *status)
		argIdx++
	}
	if keyword != nil && *keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(p.title ILIKE $%d OR p.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+*keyword+"%")
		argIdx++
	}
	if minBudget != nil {
		conditions = append(conditions, fmt.Sprintf("p.budget_max >= $%d", argIdx))
		args = append(args, *minBudget)
		argIdx++
	}
	if maxBudget != nil {
		conditions = append(conditions, fmt.Sprintf("p.budget_min <= $%d", argIdx))
		args = append(args, *maxBudget)
		argIdx++
	}
	if len(techStack) > 0 {
		conditions = append(conditions, fmt.Sprintf("p.tech_stack @> $%d::jsonb", argIdx))
		args = append(args, techStack)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects p WHERE %s", where)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}

	if total == 0 {
		return []*model.Project{}, 0, nil
	}

	// Data query with joins
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE %s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`,
		projectColumns, projectJoin, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("search projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*model.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, total, nil
}

func (r *ProjectRepository) ListByClient(ctx context.Context, clientID string, status *string, page, pageSize int) ([]*model.Project, int64, error) {
	conditions := []string{"p.client_id = $1::uuid"}
	args := []interface{}{clientID}
	argIdx := 2

	if status != nil && *status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, *status)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects p WHERE %s", where)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count client projects: %w", err)
	}

	if total == 0 {
		return []*model.Project{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE %s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`,
		projectColumns, projectJoin, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list client projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*model.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, total, nil
}

func (r *ProjectRepository) ListByDeveloper(ctx context.Context, developerID string, status *string, page, pageSize int) ([]*model.Project, int64, error) {
	// Projects where developer has bid or is assigned
	conditions := []string{
		`(p.id IN (SELECT b.project_id FROM bids b WHERE b.developer_id = $1::uuid) OR p.assigned_developer_id = $1::uuid)`,
	}
	args := []interface{}{developerID}
	argIdx := 2

	if status != nil && *status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, *status)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects p WHERE %s", where)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count developer projects: %w", err)
	}

	if total == 0 {
		return []*model.Project{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE %s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`,
		projectColumns, projectJoin, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list developer projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*model.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, total, nil
}

func (r *ProjectRepository) ListFeatured(ctx context.Context, page, pageSize int) ([]*model.Project, int64, error) {
	countQuery := `SELECT COUNT(*) FROM projects p WHERE p.featured = true AND p.status IN ('published', 'bidding', 'in_progress')`
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count featured projects: %w", err)
	}

	if total == 0 {
		return []*model.Project{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE p.featured = true AND p.status IN ('published', 'bidding', 'in_progress') ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`,
		projectColumns, projectJoin,
	)

	rows, err := r.db.Pool.Query(ctx, dataQuery, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list featured projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*model.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, total, nil
}
