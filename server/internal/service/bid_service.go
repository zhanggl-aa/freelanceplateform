package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type BidService struct {
	bidRepo      *repository.BidRepository
	projectRepo  *repository.ProjectRepository
	contractRepo *repository.ContractRepository
}

func NewBidService(
	bidRepo *repository.BidRepository,
	projectRepo *repository.ProjectRepository,
	contractRepo *repository.ContractRepository,
) *BidService {
	return &BidService{
		bidRepo:      bidRepo,
		projectRepo:  projectRepo,
		contractRepo: contractRepo,
	}
}

func (s *BidService) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	if bid.ProjectID == "" {
		return nil, errors.New("project_id is required")
	}
	if bid.DeveloperID == "" {
		return nil, errors.New("developer_id is required")
	}
	if bid.CoverLetter == "" {
		return nil, errors.New("cover_letter is required")
	}
	if bid.ProposedBudget <= 0 {
		return nil, errors.New("proposed_budget must be positive")
	}
	if bid.EstimatedDays <= 0 {
		return nil, errors.New("estimated_days must be positive")
	}
	if bid.BudgetType == "" {
		bid.BudgetType = "fixed"
	}

	validBudgetTypes := map[string]bool{"fixed": true, "hourly": true}
	if !validBudgetTypes[bid.BudgetType] {
		return nil, errors.New("budget_type must be 'fixed' or 'hourly'")
	}

	project, err := s.projectRepo.GetByID(ctx, bid.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.Status != "published" && project.Status != "bidding" && project.Status != "open" {
		return nil, errors.New("project is not open for bidding")
	}

	if project.BidDeadline != nil && time.Now().After(*project.BidDeadline) {
		return nil, errors.New("bid deadline has passed")
	}

	if project.ClientID == bid.DeveloperID {
		return nil, errors.New("cannot bid on your own project")
	}

	existingBid, err := s.bidRepo.GetByProjectAndDeveloper(ctx, bid.ProjectID, bid.DeveloperID)
	if err != nil {
		return nil, fmt.Errorf("check existing bid: %w", err)
	}
	if existingBid != nil {
		return nil, errors.New("you have already bid on this project")
	}

	created, err := s.bidRepo.Create(ctx, bid.ProjectID, bid.DeveloperID, bid.CoverLetter, bid.EstimatedDays, bid.ProposedBudget, bid.BudgetType, bid.MilestonePlan)
	if err != nil {
		return nil, fmt.Errorf("create bid: %w", err)
	}

	_ = s.projectRepo.IncrementBidCount(ctx, bid.ProjectID)

	return created, nil
}

func (s *BidService) GetByID(ctx context.Context, id string) (*model.Bid, error) {
	if id == "" {
		return nil, errors.New("bid id is required")
	}

	bid, err := s.bidRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get bid: %w", err)
	}
	return bid, nil
}

func (s *BidService) ListByProject(ctx context.Context, projectID string, page, pageSize int) ([]*model.Bid, int64, error) {
	if projectID == "" {
		return nil, 0, errors.New("project_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.bidRepo.GetByProject(ctx, projectID, page, pageSize)
}

func (s *BidService) ListByDeveloper(ctx context.Context, developerID string, page, pageSize int) ([]*model.Bid, int64, error) {
	if developerID == "" {
		return nil, 0, errors.New("developer_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.bidRepo.GetByDeveloper(ctx, developerID, nil, page, pageSize)
}

func (s *BidService) Update(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	if bid.ID == "" {
		return nil, errors.New("bid id is required")
	}

	existing, err := s.bidRepo.GetByID(ctx, bid.ID)
	if err != nil {
		return nil, fmt.Errorf("get bid: %w", err)
	}
	if existing == nil {
		return nil, errors.New("bid not found")
	}
	if existing.Status != "pending" {
		return nil, errors.New("only pending bids can be updated")
	}

	fields := map[string]interface{}{
		"cover_letter":    bid.CoverLetter,
		"estimated_days":  bid.EstimatedDays,
		"proposed_budget": bid.ProposedBudget,
		"budget_type":     bid.BudgetType,
	}
	if bid.MilestonePlan != nil {
		fields["milestone_plan"] = *bid.MilestonePlan
	}

	if err := s.bidRepo.Update(ctx, bid.ID, fields); err != nil {
		return nil, fmt.Errorf("update bid: %w", err)
	}

	return s.bidRepo.GetByID(ctx, bid.ID)
}

func (s *BidService) Withdraw(ctx context.Context, bidID string) error {
	if bidID == "" {
		return errors.New("bid id is required")
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("get bid: %w", err)
	}
	if bid == nil {
		return errors.New("bid not found")
	}
	if bid.Status != "pending" && bid.Status != "shortlisted" {
		return errors.New("only pending or shortlisted bids can be withdrawn")
	}

	return s.bidRepo.UpdateStatus(ctx, bidID, "withdrawn", nil)
}

func (s *BidService) Accept(ctx context.Context, bidID string) (*model.Contract, error) {
	if bidID == "" {
		return nil, errors.New("bid id is required")
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("get bid: %w", err)
	}
	if bid == nil {
		return nil, errors.New("bid not found")
	}
	if bid.Status != "pending" && bid.Status != "shortlisted" {
		return nil, errors.New("only pending or shortlisted bids can be accepted")
	}

	project, err := s.projectRepo.GetByID(ctx, bid.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	// Step 1: Update bid status to accepted
	if err := s.bidRepo.UpdateStatus(ctx, bidID, "accepted", nil); err != nil {
		return nil, fmt.Errorf("accept bid: %w", err)
	}

	// Step 2: Update project status and assign developer
	if err := s.projectRepo.Update(ctx, bid.ProjectID, map[string]interface{}{
		"status":                "in_progress",
		"assigned_developer_id": bid.DeveloperID,
	}); err != nil {
		return nil, fmt.Errorf("update project status: %w", err)
	}

	// Step 3: Create contract
	platformFeeRate := 0.10
	platformFee := bid.ProposedBudget * platformFeeRate
	developerPayout := bid.ProposedBudget - platformFee

	contract, err := s.contractRepo.Create(ctx,
		bid.ProjectID, project.ClientID, bid.DeveloperID, bid.ID,
		bid.ProposedBudget, platformFeeRate, platformFee, developerPayout,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("create contract: %w", err)
	}

	// Step 4: Reject all other pending bids for this project
	otherBids, _, err := s.bidRepo.GetByProject(ctx, bid.ProjectID, 1, 100)
	if err == nil {
		for _, otherBid := range otherBids {
			if otherBid.ID != bidID && (otherBid.Status == "pending" || otherBid.Status == "shortlisted") {
				_ = s.bidRepo.UpdateStatus(ctx, otherBid.ID, "rejected", nil)
			}
		}
	}

	return contract, nil
}

func (s *BidService) Reject(ctx context.Context, bidID, message string) error {
	if bidID == "" {
		return errors.New("bid id is required")
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("get bid: %w", err)
	}
	if bid == nil {
		return errors.New("bid not found")
	}
	if bid.Status != "pending" && bid.Status != "shortlisted" {
		return errors.New("only pending or shortlisted bids can be rejected")
	}

	var clientMsg *string
	if message != "" {
		clientMsg = &message
	}
	return s.bidRepo.UpdateStatus(ctx, bidID, "rejected", clientMsg)
}

func (s *BidService) Shortlist(ctx context.Context, bidID string) error {
	if bidID == "" {
		return errors.New("bid id is required")
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("get bid: %w", err)
	}
	if bid == nil {
		return errors.New("bid not found")
	}
	if bid.Status != "pending" {
		return errors.New("only pending bids can be shortlisted")
	}

	return s.bidRepo.UpdateStatus(ctx, bidID, "shortlisted", nil)
}

func (s *BidService) CounterOffer(ctx context.Context, bidID string, proposedBudget float64, estimatedDays int, message string) (*model.Bid, error) {
	if bidID == "" {
		return nil, errors.New("bid id is required")
	}
	if proposedBudget <= 0 {
		return nil, errors.New("proposed_budget must be positive")
	}
	if estimatedDays <= 0 {
		return nil, errors.New("estimated_days must be positive")
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("get bid: %w", err)
	}
	if bid == nil {
		return nil, errors.New("bid not found")
	}
	if bid.Status != "pending" && bid.Status != "shortlisted" {
		return nil, errors.New("counter-offers can only be made on pending or shortlisted bids")
	}

	fields := map[string]interface{}{
		"proposed_budget": proposedBudget,
		"estimated_days":  estimatedDays,
		"status":          "countered",
	}
	if message != "" {
		fields["client_message"] = message
	}

	if err := s.bidRepo.Update(ctx, bidID, fields); err != nil {
		return nil, fmt.Errorf("counter-offer bid: %w", err)
	}

	return s.bidRepo.GetByID(ctx, bidID)
}
