package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type MilestoneService struct {
	milestoneRepo *repository.MilestoneRepository
	paymentRepo   *repository.PaymentRepository
	walletRepo    *repository.WalletRepository
	contractRepo  *repository.ContractRepository
}

func NewMilestoneService(
	milestoneRepo *repository.MilestoneRepository,
	paymentRepo *repository.PaymentRepository,
	walletRepo *repository.WalletRepository,
	contractRepo *repository.ContractRepository,
) *MilestoneService {
	return &MilestoneService{
		milestoneRepo: milestoneRepo,
		paymentRepo:   paymentRepo,
		walletRepo:    walletRepo,
		contractRepo:  contractRepo,
	}
}

// Create adds a new milestone to a project.
func (s *MilestoneService) Create(ctx context.Context, milestone *model.ProjectMilestone) (*model.ProjectMilestone, error) {
	if milestone.ProjectID == "" {
		return nil, errors.New("project_id is required")
	}
	if milestone.Title == "" {
		return nil, errors.New("title is required")
	}
	if milestone.Amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	var deadline *string
	if milestone.Deadline != nil {
		d := milestone.Deadline.Format("2006-01-02")
		deadline = &d
	}

	sortOrder := milestone.SortOrder
	if sortOrder == 0 {
		sortOrder = 0
	}

	created, err := s.milestoneRepo.Create(ctx, milestone.ProjectID, milestone.Title, milestone.Description, milestone.Amount, deadline, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("create milestone: %w", err)
	}

	return created, nil
}

// GetByID returns a milestone by its ID.
func (s *MilestoneService) GetByID(ctx context.Context, id string) (*model.ProjectMilestone, error) {
	if id == "" {
		return nil, errors.New("milestone id is required")
	}

	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	return milestone, nil
}

// ListByProject returns all milestones for a given project.
func (s *MilestoneService) ListByProject(ctx context.Context, projectID string) ([]*model.ProjectMilestone, error) {
	if projectID == "" {
		return nil, errors.New("project_id is required")
	}

	milestones, err := s.milestoneRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list milestones: %w", err)
	}

	return milestones, nil
}

// Update modifies an existing milestone. Only pending milestones can be updated.
func (s *MilestoneService) Update(ctx context.Context, milestone *model.ProjectMilestone) (*model.ProjectMilestone, error) {
	if milestone.ID == "" {
		return nil, errors.New("milestone id is required")
	}

	existing, err := s.milestoneRepo.GetByID(ctx, milestone.ID)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	if existing == nil {
		return nil, errors.New("milestone not found")
	}
	if existing.Status != "pending" {
		return nil, errors.New("only pending milestones can be updated")
	}

	fields := map[string]interface{}{}
	if milestone.Title != "" {
		fields["title"] = milestone.Title
	}
	if milestone.Description != nil {
		fields["description"] = *milestone.Description
	}
	if milestone.Amount != 0 {
		fields["amount"] = milestone.Amount
	}
	if milestone.Deadline != nil {
		fields["deadline"] = milestone.Deadline
	}
	if milestone.SortOrder != 0 {
		fields["sort_order"] = milestone.SortOrder
	}
	if milestone.DeliverableURLs != nil {
		fields["deliverable_urls"] = milestone.DeliverableURLs
	}

	err = s.milestoneRepo.Update(ctx, milestone.ID, fields)
	if err != nil {
		return nil, fmt.Errorf("update milestone: %w", err)
	}

	updated, err := s.milestoneRepo.GetByID(ctx, milestone.ID)
	if err != nil {
		return nil, fmt.Errorf("get updated milestone: %w", err)
	}
	return updated, nil
}

// Delete removes a milestone. Only pending milestones can be deleted.
func (s *MilestoneService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("milestone id is required")
	}

	existing, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get milestone: %w", err)
	}
	if existing == nil {
		return errors.New("milestone not found")
	}
	if existing.Status != "pending" {
		return errors.New("only pending milestones can be deleted")
	}

	return s.milestoneRepo.Delete(ctx, id)
}

// Submit marks a milestone as submitted by the developer with deliverable URLs.
func (s *MilestoneService) Submit(ctx context.Context, id string, deliverableURLs []string) (*model.ProjectMilestone, error) {
	if id == "" {
		return nil, errors.New("milestone id is required")
	}

	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	if milestone == nil {
		return nil, errors.New("milestone not found")
	}
	if milestone.Status != "pending" && milestone.Status != "revision_requested" {
		return nil, errors.New("only pending or revision_requested milestones can be submitted")
	}

	now := time.Now()
	fields := map[string]interface{}{
		"status":       "submitted",
		"submitted_at": now,
	}
	if len(deliverableURLs) > 0 {
		fields["deliverable_urls"] = deliverableURLs
	}

	err = s.milestoneRepo.Update(ctx, id, fields)
	if err != nil {
		return nil, fmt.Errorf("submit milestone: %w", err)
	}

	updated, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated milestone: %w", err)
	}
	return updated, nil
}

// Approve approves a milestone and triggers the payment release workflow:
//  1. Update milestone status to 'approved'
//  2. Find the payment for this milestone
//  3. Update payment status to 'released'
//  4. Unfreeze from client wallet, add to developer wallet
//  5. Update contract released_amount
//  6. Create wallet transactions for both parties
func (s *MilestoneService) Approve(ctx context.Context, id string, feedback string) (*model.ProjectMilestone, error) {
	if id == "" {
		return nil, errors.New("milestone id is required")
	}

	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	if milestone == nil {
		return nil, errors.New("milestone not found")
	}
	if milestone.Status != "submitted" {
		return nil, errors.New("only submitted milestones can be approved")
	}

	// Step 1: Update milestone status
	now := time.Now()
	fields := map[string]interface{}{
		"status":      "approved",
		"approved_at": now,
	}
	if feedback != "" {
		fields["client_feedback"] = feedback
	}

	err = s.milestoneRepo.Update(ctx, id, fields)
	if err != nil {
		return nil, fmt.Errorf("approve milestone: %w", err)
	}

	// Get the updated milestone for returning
	updatedMilestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated milestone: %w", err)
	}

	// Step 2: Find payments for this milestone via the contract
	// Since there is no GetByMilestoneID, we use ListByContract and find the matching payment
	contract, err := s.contractRepo.GetByProjectID(ctx, milestone.ProjectID)
	if err != nil {
		return updatedMilestone, fmt.Errorf("get contract for milestone project: %w", err)
	}
	if contract == nil {
		return updatedMilestone, nil
	}

	payments, _, err := s.paymentRepo.ListByContract(ctx, contract.ID, 1, 100)
	if err != nil {
		return updatedMilestone, fmt.Errorf("list payments for contract: %w", err)
	}

	var payment *model.Payment
	for _, p := range payments {
		if p.MilestoneID != nil && *p.MilestoneID == id {
			payment = p
			break
		}
	}

	if payment != nil && payment.Status == "escrow" {
		// Step 3: Update payment status to released
		err = s.paymentRepo.UpdateStatus(ctx, payment.ID, "released")
		if err != nil {
			return updatedMilestone, fmt.Errorf("update payment status: %w", err)
		}

		// Step 4: Unfreeze from client wallet and add to developer wallet
		clientWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayerID)
		if err != nil {
			return updatedMilestone, fmt.Errorf("get client wallet: %w", err)
		}
		if clientWallet != nil {
			// Unfreeze the amount from client wallet
			err = s.walletRepo.UnfreezeAmount(ctx, clientWallet.ID, payment.Amount)
			if err != nil {
				return updatedMilestone, fmt.Errorf("unfreeze client wallet: %w", err)
			}

			// Refetch wallet to get updated balance for transaction
			clientWallet, err = s.walletRepo.GetByUserID(ctx, payment.PayerID)
			if err != nil {
				return updatedMilestone, fmt.Errorf("get updated client wallet: %w", err)
			}

			// Create wallet transaction for client (debit from frozen)
			_, err = s.walletRepo.CreateTransaction(ctx, clientWallet.ID, payment.ID, "milestone_payment", -payment.Amount, clientWallet.Balance, strPtr(fmt.Sprintf("Payment released for milestone: %s", milestone.Title)))
			if err != nil {
				return updatedMilestone, fmt.Errorf("create client wallet transaction: %w", err)
			}
		}

		developerWallet, err := s.walletRepo.GetByUserID(ctx, payment.PayeeID)
		if err != nil {
			return updatedMilestone, fmt.Errorf("get developer wallet: %w", err)
		}
		if developerWallet != nil {
			// Add net amount to developer wallet via AddDeposit
			err = s.walletRepo.AddDeposit(ctx, developerWallet.ID, payment.NetAmount)
			if err != nil {
				return updatedMilestone, fmt.Errorf("credit developer wallet: %w", err)
			}

			// Refetch wallet to get updated balance for transaction
			developerWallet, err = s.walletRepo.GetByUserID(ctx, payment.PayeeID)
			if err != nil {
				return updatedMilestone, fmt.Errorf("get updated developer wallet: %w", err)
			}

			// Create wallet transaction for developer (credit)
			_, err = s.walletRepo.CreateTransaction(ctx, developerWallet.ID, payment.ID, "milestone_payout", payment.NetAmount, developerWallet.Balance, strPtr(fmt.Sprintf("Payout for milestone: %s", milestone.Title)))
			if err != nil {
				return updatedMilestone, fmt.Errorf("create developer wallet transaction: %w", err)
			}
		}

		// Step 5: Update contract released_amount
		err = s.contractRepo.UpdateReleasedAmount(ctx, contract.ID, payment.Amount)
		if err != nil {
			return updatedMilestone, fmt.Errorf("update contract released amount: %w", err)
		}
	}

	return updatedMilestone, nil
}

// Reject sends a submitted milestone back for revision.
func (s *MilestoneService) Reject(ctx context.Context, id string, feedback string) (*model.ProjectMilestone, error) {
	if id == "" {
		return nil, errors.New("milestone id is required")
	}
	if feedback == "" {
		return nil, errors.New("feedback is required when rejecting a milestone")
	}

	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	if milestone == nil {
		return nil, errors.New("milestone not found")
	}
	if milestone.Status != "submitted" {
		return nil, errors.New("only submitted milestones can be rejected")
	}

	err = s.milestoneRepo.UpdateStatus(ctx, id, "revision_requested", &feedback)
	if err != nil {
		return nil, fmt.Errorf("reject milestone: %w", err)
	}

	updated, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated milestone: %w", err)
	}
	return updated, nil
}

// Dispute transitions a milestone to disputed status.
func (s *MilestoneService) Dispute(ctx context.Context, id string, reason string) (*model.ProjectMilestone, error) {
	if id == "" {
		return nil, errors.New("milestone id is required")
	}
	if reason == "" {
		return nil, errors.New("dispute reason is required")
	}

	milestone, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get milestone: %w", err)
	}
	if milestone == nil {
		return nil, errors.New("milestone not found")
	}

	if milestone.Status != "submitted" && milestone.Status != "approved" {
		return nil, errors.New("only submitted or approved milestones can be disputed")
	}

	err = s.milestoneRepo.UpdateStatus(ctx, id, "disputed", &reason)
	if err != nil {
		return nil, fmt.Errorf("dispute milestone: %w", err)
	}

	updated, err := s.milestoneRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated milestone: %w", err)
	}
	return updated, nil
}

// strPtr is a helper to get a pointer to a string.
func strPtr(s string) *string {
	return &s
}
