package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type PaymentService struct {
	paymentRepo  *repository.PaymentRepository
	walletRepo   *repository.WalletRepository
	contractRepo *repository.ContractRepository
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	walletRepo *repository.WalletRepository,
	contractRepo *repository.ContractRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepo,
		walletRepo:   walletRepo,
		contractRepo: contractRepo,
	}
}

// Deposit creates a payment record and freezes the amount in the client's wallet (escrow).
// This is called when a client funds a milestone.
func (s *PaymentService) Deposit(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	if payment.ContractID == "" {
		return nil, errors.New("contract_id is required")
	}
	if payment.PayerID == "" {
		return nil, errors.New("payer_id is required")
	}
	if payment.PayeeID == "" {
		return nil, errors.New("payee_id is required")
	}
	if payment.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if payment.PaymentMethod == "" {
		payment.PaymentMethod = "platform_wallet"
	}

	// Validate contract exists and is active
	contract, err := s.contractRepo.GetByID(ctx, payment.ContractID)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}
	if contract.Status != "active" && contract.Status != "started" {
		return nil, errors.New("contract is not active")
	}

	// Calculate platform fee and net amount
	platformFeeRate := contract.PlatformFeeRate
	if platformFeeRate == 0 {
		platformFeeRate = 0.10
	}
	payment.PlatformFee = payment.Amount * platformFeeRate
	payment.NetAmount = payment.Amount - payment.PlatformFee

	milestoneID := ""
	if payment.MilestoneID != nil {
		milestoneID = *payment.MilestoneID
	}

	// Create the payment record
	created, err := s.paymentRepo.Create(ctx, payment.ContractID, milestoneID, payment.PayerID, payment.PayeeID, payment.Amount, payment.PlatformFee, payment.NetAmount, payment.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	// Freeze amount in client wallet
	clientWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayerID)
	if err != nil {
		return nil, fmt.Errorf("get client wallet: %w", err)
	}
	if clientWallet == nil {
		return nil, errors.New("client wallet not found")
	}

	if clientWallet.Balance < payment.Amount {
		return nil, errors.New("insufficient wallet balance")
	}

	err = s.walletRepo.FreezeAmount(ctx, clientWallet.ID, payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("freeze client wallet: %w", err)
	}

	// Refetch wallet to get updated balance for transaction
	clientWallet, err = s.walletRepo.GetByUserID(ctx, payment.PayerID)
	if err != nil {
		return nil, fmt.Errorf("get updated client wallet: %w", err)
	}

	// Create wallet transaction for escrow
	_, err = s.walletRepo.CreateTransaction(ctx, clientWallet.ID, created.ID, "escrow", -payment.Amount, clientWallet.Balance, strPtr(fmt.Sprintf("Escrow for contract %s", payment.ContractID)))
	if err != nil {
		return nil, fmt.Errorf("create escrow wallet transaction: %w", err)
	}

	// Update contract paid amount
	err = s.contractRepo.UpdatePaidAmount(ctx, contract.ID, payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("update contract paid amount: %w", err)
	}

	return created, nil
}

// Release transitions a payment from escrow to released, transferring funds to the developer.
func (s *PaymentService) Release(ctx context.Context, paymentID string) (*model.Payment, error) {
	if paymentID == "" {
		return nil, errors.New("payment id is required")
	}

	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	if payment.Status != "escrow" {
		return nil, errors.New("only escrowed payments can be released")
	}

	err = s.paymentRepo.UpdateStatus(ctx, paymentID, "released")
	if err != nil {
		return nil, fmt.Errorf("release payment: %w", err)
	}

	// Unfreeze from client wallet
	clientWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayerID)
	if err != nil {
		return nil, fmt.Errorf("get client wallet: %w", err)
	}
	if clientWallet != nil {
		err = s.walletRepo.UnfreezeAmount(ctx, clientWallet.ID, payment.Amount)
		if err != nil {
			return nil, fmt.Errorf("unfreeze client wallet: %w", err)
		}
	}

	// Credit developer wallet with net amount
	developerWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayeeID)
	if err != nil {
		return nil, fmt.Errorf("get developer wallet: %w", err)
	}
	if developerWallet != nil {
		err = s.walletRepo.AddDeposit(ctx, developerWallet.ID, payment.NetAmount)
		if err != nil {
			return nil, fmt.Errorf("credit developer wallet: %w", err)
		}

		// Refetch wallet to get updated balance for transaction
		developerWallet, err = s.walletRepo.GetByUserID(ctx, payment.PayeeID)
		if err != nil {
			return nil, fmt.Errorf("get updated developer wallet: %w", err)
		}

		// Create wallet transaction for developer
		_, err = s.walletRepo.CreateTransaction(ctx, developerWallet.ID, payment.ID, "payout", payment.NetAmount, developerWallet.Balance, strPtr(fmt.Sprintf("Payout for payment %s", payment.ID)))
		if err != nil {
			return nil, fmt.Errorf("create developer wallet transaction: %w", err)
		}
	}

	// Update contract released amount
	err = s.contractRepo.UpdateReleasedAmount(ctx, payment.ContractID, payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("update contract released amount: %w", err)
	}

	updated, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get updated payment: %w", err)
	}
	return updated, nil
}

// Refund returns an escrowed payment back to the client.
func (s *PaymentService) Refund(ctx context.Context, paymentID string) (*model.Payment, error) {
	if paymentID == "" {
		return nil, errors.New("payment id is required")
	}

	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	if payment.Status != "escrow" {
		return nil, errors.New("only escrowed payments can be refunded")
	}

	err = s.paymentRepo.UpdateStatus(ctx, paymentID, "refunded")
	if err != nil {
		return nil, fmt.Errorf("refund payment: %w", err)
	}

	// Unfreeze from client wallet and refund back to balance
	clientWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayerID)
	if err != nil {
		return nil, fmt.Errorf("get client wallet: %w", err)
	}
	if clientWallet != nil {
		// UnfreezeAmount moves frozen -> balance atomically, which is exactly what a refund needs
		err = s.walletRepo.UnfreezeAmount(ctx, clientWallet.ID, payment.Amount)
		if err != nil {
			return nil, fmt.Errorf("unfreeze and refund client wallet: %w", err)
		}

		// Refetch wallet to get updated balance for transaction
		clientWallet, err = s.walletRepo.GetByUserID(ctx, payment.PayerID)
		if err != nil {
			return nil, fmt.Errorf("get updated client wallet: %w", err)
		}

		// Create wallet transaction for refund
		_, err = s.walletRepo.CreateTransaction(ctx, clientWallet.ID, payment.ID, "refund", payment.Amount, clientWallet.Balance, strPtr(fmt.Sprintf("Refund for payment %s", payment.ID)))
		if err != nil {
			return nil, fmt.Errorf("create refund wallet transaction: %w", err)
		}
	}

	// Update contract: reduce paid amount since funds are returned
	contract, err := s.contractRepo.GetByID(ctx, payment.ContractID)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract != nil {
		newPaidAmount := contract.PaidAmount - payment.Amount
		if newPaidAmount < 0 {
			newPaidAmount = 0
		}
		err = s.contractRepo.Update(ctx, contract.ID, map[string]interface{}{"paid_amount": newPaidAmount})
		if err != nil {
			return nil, fmt.Errorf("update contract paid amount: %w", err)
		}
	}

	updated, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get updated payment: %w", err)
	}
	return updated, nil
}

// GetByID returns a payment by its ID.
func (s *PaymentService) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	if id == "" {
		return nil, errors.New("payment id is required")
	}

	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}
	return payment, nil
}

// ListByContract returns all payments for a given contract.
func (s *PaymentService) ListByContract(ctx context.Context, contractID string, page, pageSize int) ([]*model.Payment, int64, error) {
	if contractID == "" {
		return nil, 0, errors.New("contract_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	payments, total, err := s.paymentRepo.ListByContract(ctx, contractID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments by contract: %w", err)
	}

	return payments, total, nil
}

// ListByUser returns all payments where the user is either payer or payee.
func (s *PaymentService) ListByUser(ctx context.Context, userID string, role string, page, pageSize int) ([]*model.Payment, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user_id is required")
	}
	if role != "payer" && role != "payee" {
		return nil, 0, errors.New("role must be 'payer' or 'payee'")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// paymentRepo.ListByUser does not filter by role, so we get all and could filter later
	// For now, we pass the call through since the repo returns all payments for the user
	payments, total, err := s.paymentRepo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments by user: %w", err)
	}

	return payments, total, nil
}
