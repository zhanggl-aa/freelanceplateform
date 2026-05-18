package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type ContractService struct {
	contractRepo *repository.ContractRepository
}

func NewContractService(contractRepo *repository.ContractRepository) *ContractService {
	return &ContractService{contractRepo: contractRepo}
}

// GetByID returns a contract by its ID.
func (s *ContractService) GetByID(ctx context.Context, id string) (*model.Contract, error) {
	if id == "" {
		return nil, errors.New("contract id is required")
	}

	contract, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	return contract, nil
}

// GetByProjectID returns the contract for a given project.
func (s *ContractService) GetByProjectID(ctx context.Context, projectID string) (*model.Contract, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}

	contract, err := s.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get contract by project: %w", err)
	}
	return contract, nil
}

// ListByUser returns all contracts for a user (as either client or developer).
func (s *ContractService) ListByUser(ctx context.Context, userID string, role string, page, pageSize int) ([]*model.Contract, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user id is required")
	}
	if role != "client" && role != "developer" {
		return nil, 0, errors.New("role must be 'client' or 'developer'")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	contracts, total, err := s.contractRepo.ListByUser(ctx, userID, nil, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list contracts by user: %w", err)
	}

	return contracts, total, nil
}

// Start transitions a contract from 'active' to 'started', indicating work has begun.
func (s *ContractService) Start(ctx context.Context, id string) (*model.Contract, error) {
	if id == "" {
		return nil, errors.New("contract id is required")
	}

	contract, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status != "active" {
		return nil, errors.New("only active contracts can be started")
	}

	err = s.contractRepo.UpdateStatus(ctx, id, "started")
	if err != nil {
		return nil, fmt.Errorf("start contract: %w", err)
	}

	updated, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated contract: %w", err)
	}
	return updated, nil
}

// Cancel transitions a contract to 'cancelled'. Only active or started contracts can be cancelled.
func (s *ContractService) Cancel(ctx context.Context, id string, reason string) (*model.Contract, error) {
	if id == "" {
		return nil, errors.New("contract id is required")
	}

	contract, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status != "active" && contract.Status != "started" {
		return nil, errors.New("only active or started contracts can be cancelled")
	}

	fields := map[string]interface{}{
		"status": "cancelled",
	}
	if reason != "" {
		fields["terms"] = reason
	}

	err = s.contractRepo.Update(ctx, id, fields)
	if err != nil {
		return nil, fmt.Errorf("cancel contract: %w", err)
	}

	updated, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated contract: %w", err)
	}
	return updated, nil
}

// Dispute transitions a contract to 'disputed'.
func (s *ContractService) Dispute(ctx context.Context, id string, reason string) (*model.Contract, error) {
	if id == "" {
		return nil, errors.New("contract id is required")
	}
	if reason == "" {
		return nil, errors.New("dispute reason is required")
	}

	contract, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("contract not found")
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status != "active" && contract.Status != "started" {
		return nil, errors.New("only active or started contracts can be disputed")
	}

	err = s.contractRepo.Update(ctx, id, map[string]interface{}{
		"status": "disputed",
		"terms":  reason,
	})
	if err != nil {
		return nil, fmt.Errorf("dispute contract: %w", err)
	}

	updated, err := s.contractRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated contract: %w", err)
	}
	return updated, nil
}
