package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	db *DB
}

func NewPaymentRepository(db *DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, contractID, milestoneID, payerID, payeeID string, amount, platformFee, netAmount float64, paymentMethod string) (*model.Payment, error) {
	var p model.Payment
	var milestoneIDPtr *string
	if milestoneID != "" {
		milestoneIDPtr = &milestoneID
	}
	query := `
		INSERT INTO payments (contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount,
			payment_method, external_tx_id, status, escrow_at, released_at, refunded_at, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		contractID, milestoneIDPtr, payerID, payeeID, amount, platformFee, netAmount, paymentMethod,
	).Scan(
		&p.ID, &p.ContractID, &p.MilestoneID, &p.PayerID, &p.PayeeID,
		&p.Amount, &p.PlatformFee, &p.NetAmount, &p.PaymentMethod,
		&p.ExternalTxID, &p.Status, &p.EscrowAt, &p.ReleasedAt, &p.RefundedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}
	return &p, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	var p model.Payment
	query := `
		SELECT id, contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount,
			payment_method, external_tx_id, status, escrow_at, released_at, refunded_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.ContractID, &p.MilestoneID, &p.PayerID, &p.PayeeID,
		&p.Amount, &p.PlatformFee, &p.NetAmount, &p.PaymentMethod,
		&p.ExternalTxID, &p.Status, &p.EscrowAt, &p.ReleasedAt, &p.RefundedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get payment by id: %w", err)
	}
	return &p, nil
}

func (r *PaymentRepository) ListByContract(ctx context.Context, contractID string, page, pageSize int) ([]*model.Payment, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments WHERE contract_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, contractID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count payments by contract: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount,
			payment_method, external_tx_id, status, escrow_at, released_at, refunded_at, created_at, updated_at
		FROM payments
		WHERE contract_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, contractID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments by contract: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.ID, &p.ContractID, &p.MilestoneID, &p.PayerID, &p.PayeeID,
			&p.Amount, &p.PlatformFee, &p.NetAmount, &p.PaymentMethod,
			&p.ExternalTxID, &p.Status, &p.EscrowAt, &p.ReleasedAt, &p.RefundedAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate payments: %w", err)
	}
	return payments, total, nil
}

func (r *PaymentRepository) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Payment, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments WHERE payer_id = $1 OR payee_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count payments by user: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount,
			payment_method, external_tx_id, status, escrow_at, released_at, refunded_at, created_at, updated_at
		FROM payments
		WHERE payer_id = $1 OR payee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments by user: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.ID, &p.ContractID, &p.MilestoneID, &p.PayerID, &p.PayeeID,
			&p.Amount, &p.PlatformFee, &p.NetAmount, &p.PaymentMethod,
			&p.ExternalTxID, &p.Status, &p.EscrowAt, &p.ReleasedAt, &p.RefundedAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate payments: %w", err)
	}
	return payments, total, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE payments SET status = $2 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

func (r *PaymentRepository) UpdateExternalTxID(ctx context.Context, id string, externalTxID string) error {
	query := `UPDATE payments SET external_tx_id = $2 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, externalTxID)
	if err != nil {
		return fmt.Errorf("update payment external tx id: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

func (r *PaymentRepository) ListAll(ctx context.Context, page, pageSize int) ([]*model.Payment, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments`
	if err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count all payments: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, contract_id, milestone_id, payer_id, payee_id, amount, platform_fee, net_amount,
			payment_method, external_tx_id, status, escrow_at, released_at, refunded_at, created_at, updated_at
		FROM payments
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list all payments: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.ID, &p.ContractID, &p.MilestoneID, &p.PayerID, &p.PayeeID,
			&p.Amount, &p.PlatformFee, &p.NetAmount, &p.PaymentMethod,
			&p.ExternalTxID, &p.Status, &p.EscrowAt, &p.ReleasedAt, &p.RefundedAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate payments: %w", err)
	}
	return payments, total, nil
}

func (r *PaymentRepository) SumByStatus(ctx context.Context, status string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = $1`
	err := r.db.Pool.QueryRow(ctx, query, status).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum payments by status: %w", err)
	}
	return total, nil
}

func (r *PaymentRepository) SumPlatformFees(ctx context.Context) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(platform_fee), 0) FROM payments WHERE status IN ('released', 'escrow')`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum platform fees: %w", err)
	}
	return total, nil
}

func (r *PaymentRepository) SumEscrowByContract(ctx context.Context, contractID string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE contract_id = $1 AND status = 'escrow'`
	err := r.db.Pool.QueryRow(ctx, query, contractID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum escrow by contract: %w", err)
	}
	return total, nil
}
