package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type WalletRepository struct {
	db *DB
}

func NewWalletRepository(db *DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetByUserID(ctx context.Context, userID string) (*model.PlatformWallet, error) {
	var w model.PlatformWallet
	query := `
		SELECT id, user_id, balance, frozen_amount, total_deposited, total_withdrawn, created_at, updated_at
		FROM platform_wallets
		WHERE user_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.FrozenAmount,
		&w.TotalDeposited, &w.TotalWithdrawn, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet by user id: %w", err)
	}
	return &w, nil
}

func (r *WalletRepository) Create(ctx context.Context, userID string) (*model.PlatformWallet, error) {
	var w model.PlatformWallet
	query := `
		INSERT INTO platform_wallets (user_id)
		VALUES ($1)
		RETURNING id, user_id, balance, frozen_amount, total_deposited, total_withdrawn, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.FrozenAmount,
		&w.TotalDeposited, &w.TotalWithdrawn, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}
	return &w, nil
}

func (r *WalletRepository) GetOrCreateByUserID(ctx context.Context, userID string) (*model.PlatformWallet, error) {
	w, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return w, nil
	}
	return r.Create(ctx, userID)
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID string, balance, frozenAmount float64) error {
	query := `UPDATE platform_wallets SET balance = $2, frozen_amount = $3 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, walletID, balance, frozenAmount)
	if err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found")
	}
	return nil
}

func (r *WalletRepository) AddDeposit(ctx context.Context, walletID string, amount float64) error {
	query := `
		UPDATE platform_wallets
		SET balance = balance + $2, total_deposited = total_deposited + $2
		WHERE id = $1
	`
	tag, err := r.db.Pool.Exec(ctx, query, walletID, amount)
	if err != nil {
		return fmt.Errorf("add deposit: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found")
	}
	return nil
}

func (r *WalletRepository) AddWithdrawal(ctx context.Context, walletID string, amount float64) error {
	query := `
		UPDATE platform_wallets
		SET balance = balance - $2, total_withdrawn = total_withdrawn + $2
		WHERE id = $1 AND balance >= $2
	`
	tag, err := r.db.Pool.Exec(ctx, query, walletID, amount)
	if err != nil {
		return fmt.Errorf("add withdrawal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found or insufficient balance")
	}
	return nil
}

func (r *WalletRepository) FreezeAmount(ctx context.Context, walletID string, amount float64) error {
	query := `
		UPDATE platform_wallets
		SET balance = balance - $2, frozen_amount = frozen_amount + $2
		WHERE id = $1 AND balance >= $2
	`
	tag, err := r.db.Pool.Exec(ctx, query, walletID, amount)
	if err != nil {
		return fmt.Errorf("freeze amount: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found or insufficient balance")
	}
	return nil
}

func (r *WalletRepository) UnfreezeAmount(ctx context.Context, walletID string, amount float64) error {
	query := `
		UPDATE platform_wallets
		SET frozen_amount = frozen_amount - $2, balance = balance + $2
		WHERE id = $1 AND frozen_amount >= $2
	`
	tag, err := r.db.Pool.Exec(ctx, query, walletID, amount)
	if err != nil {
		return fmt.Errorf("unfreeze amount: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found or insufficient frozen amount")
	}
	return nil
}

func (r *WalletRepository) CreateTransaction(ctx context.Context, walletID, paymentID, txType string, amount, balanceAfter float64, description *string) (*model.WalletTransaction, error) {
	var t model.WalletTransaction
	var paymentIDPtr *string
	if paymentID != "" {
		paymentIDPtr = &paymentID
	}
	query := `
		INSERT INTO wallet_transactions (wallet_id, payment_id, type, amount, balance_after, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, wallet_id, payment_id, type, amount, balance_after, description, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		walletID, paymentIDPtr, txType, amount, balanceAfter, description,
	).Scan(
		&t.ID, &t.WalletID, &t.PaymentID, &t.Type, &t.Amount, &t.BalanceAfter, &t.Description, &t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create wallet transaction: %w", err)
	}
	return &t, nil
}

func (r *WalletRepository) ListTransactions(ctx context.Context, walletID string, page, pageSize int) ([]*model.WalletTransaction, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM wallet_transactions WHERE wallet_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, walletID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count wallet transactions: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, wallet_id, payment_id, type, amount, balance_after, description, created_at
		FROM wallet_transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, walletID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list wallet transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.WalletTransaction
	for rows.Next() {
		var t model.WalletTransaction
		if err := rows.Scan(
			&t.ID, &t.WalletID, &t.PaymentID, &t.Type, &t.Amount, &t.BalanceAfter, &t.Description, &t.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan wallet transaction: %w", err)
		}
		transactions = append(transactions, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate wallet transactions: %w", err)
	}
	return transactions, total, nil
}
