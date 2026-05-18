package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type MilestoneRepository struct {
	db *DB
}

func NewMilestoneRepository(db *DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

const milestoneColumns = `
	m.id::text, m.project_id::text, m.title, m.description,
	m.amount, m.deadline, m.status, m.sort_order,
	m.deliverable_urls, m.client_feedback, m.submitted_at, m.approved_at,
	m.created_at, m.updated_at
`

func scanMilestone(row pgx.Row) (*model.ProjectMilestone, error) {
	ms := &model.ProjectMilestone{}
	err := row.Scan(
		&ms.ID, &ms.ProjectID, &ms.Title, &ms.Description,
		&ms.Amount, &ms.Deadline, &ms.Status, &ms.SortOrder,
		&ms.DeliverableURLs, &ms.ClientFeedback, &ms.SubmittedAt, &ms.ApprovedAt,
		&ms.CreatedAt, &ms.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (r *MilestoneRepository) Create(ctx context.Context, projectID, title string, description *string, amount float64, deadline *string, sortOrder int) (*model.ProjectMilestone, error) {
	query := `
		INSERT INTO project_milestones (project_id, title, description, amount, deadline, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, project_id::text, title, description,
			amount, deadline, status, sort_order,
			deliverable_urls, client_feedback, submitted_at, approved_at,
			created_at, updated_at
	`
	ms := &model.ProjectMilestone{}
	err := r.db.Pool.QueryRow(ctx, query,
		projectID, title, description, amount, deadline, sortOrder,
	).Scan(
		&ms.ID, &ms.ProjectID, &ms.Title, &ms.Description,
		&ms.Amount, &ms.Deadline, &ms.Status, &ms.SortOrder,
		&ms.DeliverableURLs, &ms.ClientFeedback, &ms.SubmittedAt, &ms.ApprovedAt,
		&ms.CreatedAt, &ms.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create milestone: %w", err)
	}
	return ms, nil
}

func (r *MilestoneRepository) GetByID(ctx context.Context, id string) (*model.ProjectMilestone, error) {
	query := `SELECT ` + milestoneColumns + ` FROM project_milestones m WHERE m.id = $1`
	ms, err := scanMilestone(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get milestone by id: %w", err)
	}
	return ms, nil
}

func (r *MilestoneRepository) ListByProject(ctx context.Context, projectID string) ([]*model.ProjectMilestone, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM project_milestones m WHERE m.project_id = $1::uuid ORDER BY m.sort_order ASC, m.created_at ASC`,
		milestoneColumns,
	)
	rows, err := r.db.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("list milestones by project: %w", err)
	}
	defer rows.Close()

	milestones := make([]*model.ProjectMilestone, 0)
	for rows.Next() {
		ms, err := scanMilestone(rows)
		if err != nil {
			return nil, fmt.Errorf("scan milestone: %w", err)
		}
		milestones = append(milestones, ms)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate milestones: %w", err)
	}

	return milestones, nil
}

func (r *MilestoneRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"title": true, "description": true, "amount": true, "deadline": true,
		"status": true, "sort_order": true, "deliverable_urls": true,
		"client_feedback": true, "submitted_at": true, "approved_at": true,
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
	query := fmt.Sprintf("UPDATE project_milestones SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update milestone: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("milestone not found: %s", id)
	}
	return nil
}

func (r *MilestoneRepository) UpdateStatus(ctx context.Context, id, status string, feedback *string) error {
	var query string
	var args []interface{}

	if feedback != nil {
		query = `UPDATE project_milestones SET status = $1, client_feedback = $2 WHERE id = $3`
		args = []interface{}{status, feedback, id}
	} else {
		query = `UPDATE project_milestones SET status = $1 WHERE id = $2`
		args = []interface{}{status, id}
	}

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update milestone status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("milestone not found: %s", id)
	}
	return nil
}

func (r *MilestoneRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM project_milestones WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete milestone: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("milestone not found: %s", id)
	}
	return nil
}

func (r *MilestoneRepository) SumApprovedAmount(ctx context.Context, projectID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) FROM project_milestones
		WHERE project_id = $1::uuid AND status = 'approved'
	`
	var total float64
	if err := r.db.Pool.QueryRow(ctx, query, projectID).Scan(&total); err != nil {
		return 0, fmt.Errorf("sum approved amount: %w", err)
	}
	return total, nil
}
