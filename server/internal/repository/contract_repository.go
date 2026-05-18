package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type ContractRepository struct {
	db *DB
}

func NewContractRepository(db *DB) *ContractRepository {
	return &ContractRepository{db: db}
}

const contractColumns = `
	c.id::text, c.project_id::text, c.client_id::text, c.developer_id::text,
	c.bid_id::text, c.total_amount, c.platform_fee_rate, c.platform_fee,
	c.developer_payout, c.paid_amount, c.released_amount, c.status,
	c.start_date, c.end_date, c.terms, c.created_at, c.updated_at
`

func scanContract(row pgx.Row) (*model.Contract, error) {
	ct := &model.Contract{}
	err := row.Scan(
		&ct.ID, &ct.ProjectID, &ct.ClientID, &ct.DeveloperID,
		&ct.BidID, &ct.TotalAmount, &ct.PlatformFeeRate, &ct.PlatformFee,
		&ct.DeveloperPayout, &ct.PaidAmount, &ct.ReleasedAmount, &ct.Status,
		&ct.StartDate, &ct.EndDate, &ct.Terms, &ct.CreatedAt, &ct.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (r *ContractRepository) Create(ctx context.Context, projectID, clientID, developerID, bidID string, totalAmount, platformFeeRate, platformFee, developerPayout float64, startDate string) (*model.Contract, error) {
	query := `
		INSERT INTO contracts (project_id, client_id, developer_id, bid_id, total_amount, platform_fee_rate, platform_fee, developer_payout, start_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text, project_id::text, client_id::text, developer_id::text,
			bid_id::text, total_amount, platform_fee_rate, platform_fee,
			developer_payout, paid_amount, released_amount, status,
			start_date, end_date, terms, created_at, updated_at
	`
	ct := &model.Contract{}
	err := r.db.Pool.QueryRow(ctx, query,
		projectID, clientID, developerID, bidID, totalAmount, platformFeeRate, platformFee, developerPayout, startDate,
	).Scan(
		&ct.ID, &ct.ProjectID, &ct.ClientID, &ct.DeveloperID,
		&ct.BidID, &ct.TotalAmount, &ct.PlatformFeeRate, &ct.PlatformFee,
		&ct.DeveloperPayout, &ct.PaidAmount, &ct.ReleasedAmount, &ct.Status,
		&ct.StartDate, &ct.EndDate, &ct.Terms, &ct.CreatedAt, &ct.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create contract: %w", err)
	}
	return ct, nil
}

func (r *ContractRepository) GetByID(ctx context.Context, id string) (*model.Contract, error) {
	query := `SELECT ` + contractColumns + ` FROM contracts c WHERE c.id = $1`
	ct, err := scanContract(r.db.Pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get contract by id: %w", err)
	}
	return ct, nil
}

func (r *ContractRepository) GetByProjectID(ctx context.Context, projectID string) (*model.Contract, error) {
	query := `SELECT ` + contractColumns + ` FROM contracts c WHERE c.project_id = $1::uuid`
	ct, err := scanContract(r.db.Pool.QueryRow(ctx, query, projectID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get contract by project id: %w", err)
	}
	return ct, nil
}

func (r *ContractRepository) ListByUser(ctx context.Context, userID string, status *string, page, pageSize int) ([]*model.Contract, int64, error) {
	conditions := []string{"(c.client_id = $1::uuid OR c.developer_id = $1::uuid)"}
	args := []interface{}{userID}
	argIdx := 2

	if status != nil && *status != "" {
		conditions = append(conditions, fmt.Sprintf("c.status = $%d", argIdx))
		args = append(args, *status)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contracts c WHERE %s", where)
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user contracts: %w", err)
	}

	if total == 0 {
		return []*model.Contract{}, 0, nil
	}

	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(
		`SELECT %s FROM contracts c WHERE %s ORDER BY c.created_at DESC LIMIT $%d OFFSET $%d`,
		contractColumns, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list user contracts: %w", err)
	}
	defer rows.Close()

	contracts := make([]*model.Contract, 0)
	for rows.Next() {
		ct, err := scanContract(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan contract: %w", err)
		}
		contracts = append(contracts, ct)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate contracts: %w", err)
	}

	return contracts, total, nil
}

func (r *ContractRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"total_amount": true, "platform_fee_rate": true, "platform_fee": true,
		"developer_payout": true, "paid_amount": true, "released_amount": true,
		"status": true, "end_date": true, "terms": true,
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
	query := fmt.Sprintf("UPDATE contracts SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update contract: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contract not found: %s", id)
	}
	return nil
}

func (r *ContractRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE contracts SET status = $1 WHERE id = $2`
	tag, err := r.db.Pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update contract status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contract not found: %s", id)
	}
	return nil
}

func (r *ContractRepository) UpdatePaidAmount(ctx context.Context, id string, amount float64) error {
	query := `UPDATE contracts SET paid_amount = paid_amount + $1 WHERE id = $2`
	tag, err := r.db.Pool.Exec(ctx, query, amount, id)
	if err != nil {
		return fmt.Errorf("update contract paid amount: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contract not found: %s", id)
	}
	return nil
}

func (r *ContractRepository) UpdateReleasedAmount(ctx context.Context, id string, amount float64) error {
	query := `UPDATE contracts SET released_amount = released_amount + $1 WHERE id = $2`
	tag, err := r.db.Pool.Exec(ctx, query, amount, id)
	if err != nil {
		return fmt.Errorf("update contract released amount: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contract not found: %s", id)
	}
	return nil
}
