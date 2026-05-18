package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

// GetBalance returns the wallet balance for a user, creating one if it does not exist.
func (s *WalletService) GetBalance(ctx context.Context, userID string) (*model.PlatformWallet, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		// Create wallet for user if it does not exist yet
		wallet, err = s.walletRepo.Create(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("create wallet: %w", err)
		}
	}

	return wallet, nil
}

// ListTransactions returns paginated wallet transactions for a user.
func (s *WalletService) ListTransactions(ctx context.Context, userID string, page, pageSize int) ([]*model.WalletTransaction, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Ensure wallet exists
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return []*model.WalletTransaction{}, 0, nil
	}

	transactions, total, err := s.walletRepo.ListTransactions(ctx, wallet.ID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list wallet transactions: %w", err)
	}

	return transactions, total, nil
}

// Withdraw transfers funds from the user's wallet balance to an external account.
// This decreases the wallet balance and creates a withdrawal transaction record.
func (s *WalletService) Withdraw(ctx context.Context, userID string, amount float64, description string) (*model.WalletTransaction, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	if wallet.Balance < amount {
		return nil, errors.New("insufficient wallet balance")
	}

	// Minimum withdrawal amount check
	if amount < 10 {
		return nil, errors.New("minimum withdrawal amount is 10")
	}

	// Use AddWithdrawal which atomically decrements balance and increments total_withdrawn
	err = s.walletRepo.AddWithdrawal(ctx, wallet.ID, amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw from wallet: %w", err)
	}

	// Refetch wallet to get updated balance for transaction
	wallet, err = s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get updated wallet: %w", err)
	}

	// Create withdrawal transaction
	desc := "Withdrawal"
	if description != "" {
		desc = description
	}

	transaction, err := s.walletRepo.CreateTransaction(ctx, wallet.ID, "", "withdrawal", -amount, wallet.Balance, &desc)
	if err != nil {
		return nil, fmt.Errorf("create withdrawal transaction: %w", err)
	}

	return transaction, nil
}
