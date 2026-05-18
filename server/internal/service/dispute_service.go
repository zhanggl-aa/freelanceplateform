package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type DisputeService struct {
	disputeRepo  *repository.DisputeRepository
	contractRepo *repository.ContractRepository
}

func NewDisputeService(
	disputeRepo *repository.DisputeRepository,
	contractRepo *repository.ContractRepository,
) *DisputeService {
	return &DisputeService{
		disputeRepo:  disputeRepo,
		contractRepo: contractRepo,
	}
}

// Create files a new dispute against a contract or milestone.
func (s *DisputeService) Create(ctx context.Context, dispute *model.Dispute) (*model.Dispute, error) {
	if dispute.ContractID == "" {
		return nil, errors.New("contract_id is required")
	}
	if dispute.ReporterID == "" {
		return nil, errors.New("reporter_id is required")
	}
	if dispute.ReportedID == "" {
		return nil, errors.New("reported_id is required")
	}
	if dispute.Reason == "" {
		return nil, errors.New("reason is required")
	}
	if dispute.ReporterID == dispute.ReportedID {
		return nil, errors.New("reporter and reported cannot be the same user")
	}

	// Validate contract exists and is active
	contract, err := s.contractRepo.GetByID(ctx, dispute.ContractID)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}
	if contract.Status != "active" && contract.Status != "started" {
		return nil, errors.New("can only dispute active or started contracts")
	}

	// Validate reporter is part of the contract
	if contract.ClientID != dispute.ReporterID && contract.DeveloperID != dispute.ReporterID {
		return nil, errors.New("reporter is not part of this contract")
	}

	// Validate reported user is the other party
	if contract.ClientID != dispute.ReportedID && contract.DeveloperID != dispute.ReportedID {
		return nil, errors.New("reported user is not part of this contract")
	}

	// Check for existing open dispute on this contract
	// Since there is no GetOpenByContract, we use ListByContract and check for open status
	existingDisputes, err := s.disputeRepo.ListByContract(ctx, dispute.ContractID)
	if err != nil {
		return nil, fmt.Errorf("check existing dispute: %w", err)
	}
	for _, d := range existingDisputes {
		if d.Status == "open" || d.Status == "under_review" {
			return nil, errors.New("there is already an open dispute on this contract")
		}
	}

	if dispute.EvidenceURLs == nil {
		dispute.EvidenceURLs = []string{}
	}

	milestoneID := ""
	if dispute.MilestoneID != nil {
		milestoneID = *dispute.MilestoneID
	}

	created, err := s.disputeRepo.Create(ctx, dispute.ContractID, milestoneID, dispute.ReporterID, dispute.ReportedID, dispute.Reason, dispute.EvidenceURLs)
	if err != nil {
		return nil, fmt.Errorf("create dispute: %w", err)
	}

	// Update contract status to disputed
	err = s.contractRepo.UpdateStatus(ctx, dispute.ContractID, "disputed")
	if err != nil {
		return nil, fmt.Errorf("update contract status: %w", err)
	}

	return created, nil
}

// GetByID returns a dispute by its ID.
func (s *DisputeService) GetByID(ctx context.Context, id string) (*model.Dispute, error) {
	if id == "" {
		return nil, errors.New("dispute id is required")
	}

	dispute, err := s.disputeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get dispute: %w", err)
	}
	return dispute, nil
}

// List returns paginated disputes, optionally filtered by status.
func (s *DisputeService) List(ctx context.Context, status string, page, pageSize int) ([]*model.Dispute, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	if status != "" {
		validStatuses := map[string]bool{"open": true, "under_review": true, "resolved": true, "dismissed": true}
		if !validStatuses[status] {
			return nil, 0, errors.New("invalid status filter")
		}
	}

	disputes, total, err := s.disputeRepo.List(ctx, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list disputes: %w", err)
	}

	return disputes, total, nil
}

// Resolve resolves a dispute (admin action). The resolution can be:
// - "favor_client": funds returned to client
// - "favor_developer": funds released to developer
// - "split": funds split between both parties
// - "dismissed": dispute dismissed
func (s *DisputeService) Resolve(ctx context.Context, disputeID, resolvedBy, resolution, resolutionType string) (*model.Dispute, error) {
	if disputeID == "" {
		return nil, errors.New("dispute id is required")
	}
	if resolvedBy == "" {
		return nil, errors.New("resolved_by is required")
	}
	if resolution == "" {
		return nil, errors.New("resolution is required")
	}

	validResolutionTypes := map[string]bool{"favor_client": true, "favor_developer": true, "split": true, "dismissed": true}
	if !validResolutionTypes[resolutionType] {
		return nil, errors.New("resolution_type must be 'favor_client', 'favor_developer', 'split', or 'dismissed'")
	}

	dispute, err := s.disputeRepo.GetByID(ctx, disputeID)
	if err != nil {
		return nil, fmt.Errorf("get dispute: %w", err)
	}
	if dispute == nil {
		return nil, errors.New("dispute not found")
	}
	if dispute.Status != "open" && dispute.Status != "under_review" {
		return nil, errors.New("only open or under_review disputes can be resolved")
	}

	// Use UpdateResolution which sets resolution, resolution_type, resolved_by, resolved_at, and status='resolved'
	err = s.disputeRepo.UpdateResolution(ctx, disputeID, resolution, resolutionType, resolvedBy)
	if err != nil {
		return nil, fmt.Errorf("resolve dispute: %w", err)
	}

	// Fetch the updated dispute
	updated, err := s.disputeRepo.GetByID(ctx, disputeID)
	if err != nil {
		return nil, fmt.Errorf("get updated dispute: %w", err)
	}

	// Update contract status based on resolution
	contract, err := s.contractRepo.GetByID(ctx, dispute.ContractID)
	if err != nil {
		return updated, fmt.Errorf("get contract: %w", err)
	}
	if contract != nil {
		var newStatus string
		switch resolutionType {
		case "dismissed":
			// Restore contract to started status
			newStatus = "started"
		case "favor_client", "favor_developer", "split":
			// Contract transitions to a resolution state
			newStatus = "dispute_resolved"
		}

		if newStatus != "" {
			err = s.contractRepo.UpdateStatus(ctx, contract.ID, newStatus)
			if err != nil {
				return updated, fmt.Errorf("update contract after dispute resolution: %w", err)
			}
		}
	}

	return updated, nil
}
