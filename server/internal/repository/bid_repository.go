package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type BidRepository struct {
	db *DB
}

func NewBidRepository(db *DB) *BidRepository {
	return &BidRepository{db: db}
}

const bidColumns = `
	b.id::text, b.project_id::text, b.developer_id::text, b.cover_letter,
	b.estimated_days, b.proposed_budget, b.budget_type, b.milestone_plan,
	b.status, b.client_message, b.created_at, b.updated_at,
	u.nickname, u.avatar_url
`

const bidJoin = `
	FROM bids b
	JOIN users u ON u.id = b.developer_id
`

func scanBid(row pgx.Row) (*model.Bid, error) {
	b := &model.Bid{}
	var developerName *string
	var developerAvatar *string

	err := row.Scan(
		&b.ID, &b.ProjectID, &b.DeveloperID, &b.CoverLetter,
		&b.EstimatedDays, &b.ProposedBudget, &b.BudgetType, &b.MilestonePlan,
		&b.Status, &b.ClientMessage, &b.CreatedAt, &b.UpdatedAt,
		&developerName, &developerAvatar,
	)
	if err != nil {
		return nil, err
	}
	b.DeveloperName = developerName
	b.DeveloperAvatar = developerAvatar
	return b, nil
}

func (r *BidRepository) Create(ctx context.Context, projectID, developerID, coverLetter string, estimatedDays int, proposedBudget float64, budgetType string, milestonePlan *string) (*model.Bid, error) {
	query := `
		INSERT INTO bids (project_id, developer_id, cover_letter, estimated_days, proposed_budget, budget_type, milestone_plan)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, project_id::text, developer_id::text, cover_letter,
			estimated_days, proposed_budget, budget_type, milestone_plan,
			status, client_message, created_at, updated_at
	`
	b := &model.Bid{}
	err := r.db.Pool.QueryRow(ctx, query,
		projectID, developerID, coverLetter, estimatedDays, proposedBudget, budgetType, milestonePlan,
	).Scan(
		&b.ID, &b.ProjectID, &b.DeveloperID, &b.CoverLetter,
		&b.EstimatedDays, &b.ProposedBudget, &b.BudgetType, &b.MilestonePlan,
		&b.Status, &b.ClientMessage, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create bid: %w", err)
	}
	return b, nil
}

func (r *BidRepository) GetByID(ctx context.Context, id string) (*model.Bid, error) {
	query := `SELECT ` + bidColumns + bidJoin + ` WHERE b.id = $1`
	b, err := scanBid(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get bid by id: %w", err)
	}
	return b, nil
}

func (r *BidRepository) GetByProject(ctx context.Context, projectID string, page, pageSize int) ([]*model.Bid, int64, error) {
	countQuery := `SELECT COUNT(*) FROM bids WHERE project_id = $1::uuid`
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, projectID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count project bids: %w", err)
	}

	if total == 0 {
		return []*model.Bid{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE b.project_id = $1::uuid ORDER BY b.created_at DESC LIMIT $2 OFFSET $3`,
		bidColumns, bidJoin,
	)

	rows, err := r.db.Pool.Query(ctx, dataQuery, projectID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get project bids: %w", err)
	}
	defer rows.Close()

	bids := make([]*model.Bid, 0)
	for rows.Next() {
		b, err := scanBid(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan bid: %w", err)
		}
		bids = append(bids, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate bids: %w", err)
	}

	return bids, total, nil
}

func (r *BidRepository) GetByDeveloper(ctx context.Context, developerID string, status *string, page, pageSize int) ([]*model.Bid, int64, error) {
	conditions := []string{"b.developer_id = $1::uuid"}
	args := []interface{}{developerID}
	argIdx := 2

	if status != nil && *status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", argIdx))
		args = append(args, *status)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM bids b WHERE %s", where)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count developer bids: %w", err)
	}

	if total == 0 {
		return []*model.Bid{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s %s WHERE %s ORDER BY b.created_at DESC LIMIT $%d OFFSET $%d`,
		bidColumns, bidJoin, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get developer bids: %w", err)
	}
	defer rows.Close()

	bids := make([]*model.Bid, 0)
	for rows.Next() {
		b, err := scanBid(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan bid: %w", err)
		}
		bids = append(bids, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate bids: %w", err)
	}

	return bids, total, nil
}

func (r *BidRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"cover_letter": true, "estimated_days": true, "proposed_budget": true,
		"budget_type": true, "milestone_plan": true, "status": true,
		"client_message": true,
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
	query := fmt.Sprintf("UPDATE bids SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update bid: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bid not found: %s", id)
	}
	return nil
}

func (r *BidRepository) UpdateStatus(ctx context.Context, id, status string, clientMessage *string) error {
	query := `UPDATE bids SET status = $1, client_message = $2 WHERE id = $3`
	tag, err := r.db.Pool.Exec(ctx, query, status, clientMessage, id)
	if err != nil {
		return fmt.Errorf("update bid status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bid not found: %s", id)
	}
	return nil
}

func (r *BidRepository) GetByProjectAndDeveloper(ctx context.Context, projectID, developerID string) (*model.Bid, error) {
	query := `SELECT ` + bidColumns + bidJoin + ` WHERE b.project_id = $1::uuid AND b.developer_id = $2::uuid`
	b, err := scanBid(r.db.Pool.QueryRow(ctx, query, projectID, developerID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get bid by project and developer: %w", err)
	}
	return b, nil
}

func (r *BidRepository) CountByProject(ctx context.Context, projectID string) (int, error) {
	query := `SELECT COUNT(*) FROM bids WHERE project_id = $1::uuid`
	var count int
	if err := r.db.Pool.QueryRow(ctx, query, projectID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count bids by project: %w", err)
	}
	return count, nil
}
