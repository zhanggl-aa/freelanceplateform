package service

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type AdminService struct {
	adminRepo   *repository.AdminRepository
	userRepo    *repository.UserRepository
	projectRepo *repository.ProjectRepository
	contractRepo *repository.ContractRepository
	paymentRepo *repository.PaymentRepository
	disputeRepo *repository.DisputeRepository
}

func NewAdminService(
	adminRepo *repository.AdminRepository,
	userRepo *repository.UserRepository,
	projectRepo *repository.ProjectRepository,
	contractRepo *repository.ContractRepository,
	paymentRepo *repository.PaymentRepository,
	disputeRepo *repository.DisputeRepository,
) *AdminService {
	return &AdminService{
		adminRepo:    adminRepo,
		userRepo:     userRepo,
		projectRepo:  projectRepo,
		contractRepo: contractRepo,
		paymentRepo:  paymentRepo,
		disputeRepo:  disputeRepo,
	}
}

func (s *AdminService) Dashboard(ctx context.Context, adminID string) (map[string]interface{}, error) {
	return s.adminRepo.DashboardStats(ctx)
}

func (s *AdminService) ListUsers(ctx context.Context, search, status *string, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	searchStr := ""
	if search != nil {
		searchStr = *search
	}
	statusStr := ""
	if status != nil {
		statusStr = *status
	}

	return s.adminRepo.ListUsers(ctx, searchStr, statusStr, page, pageSize)
}

func (s *AdminService) UpdateUserStatus(ctx context.Context, adminID, userID, status string) error {
	return s.adminRepo.UpdateUserStatus(ctx, userID, status)
}

func (s *AdminService) ListProjects(ctx context.Context, page, pageSize int) ([]*model.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.projectRepo.Search(ctx, nil, nil, nil, nil, nil, nil, page, pageSize)
}

func (s *AdminService) ModerateProject(ctx context.Context, adminID, projectID string, body map[string]interface{}) error {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	fields := make(map[string]interface{})
	if action, ok := body["action"].(string); ok {
		switch action {
		case "feature":
			fields["featured"] = true
		case "unfeature":
			fields["featured"] = false
		case "suspend":
			fields["status"] = "suspended"
		case "restore":
			if project.Status != "suspended" {
				return fmt.Errorf("project is not suspended")
			}
			fields["status"] = "published"
		default:
			return fmt.Errorf("action must be 'feature', 'unfeature', 'suspend', or 'restore'")
		}
	} else {
		// Allow direct field updates from body
		for k, v := range body {
			fields[k] = v
		}
	}

	return s.projectRepo.Update(ctx, projectID, fields)
}

func (s *AdminService) ListDisputes(ctx context.Context, page, pageSize int) ([]*model.Dispute, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.disputeRepo.List(ctx, "", page, pageSize)
}

func (s *AdminService) ResolveDispute(ctx context.Context, adminID, disputeID, resolution, resolutionType string) error {
	dispute, err := s.disputeRepo.GetByID(ctx, disputeID)
	if err != nil {
		return fmt.Errorf("get dispute: %w", err)
	}
	if dispute == nil {
		return fmt.Errorf("dispute not found")
	}
	if dispute.Status != "open" && dispute.Status != "under_review" {
		return fmt.Errorf("only open or under_review disputes can be resolved")
	}

	if err := s.disputeRepo.UpdateResolution(ctx, disputeID, resolution, resolutionType, adminID); err != nil {
		return fmt.Errorf("resolve dispute: %w", err)
	}

	contractStatus := "active"
	if resolutionType != "dismissed" {
		contractStatus = "dispute_resolved"
	}
	return s.contractRepo.UpdateStatus(ctx, dispute.ContractID, contractStatus)
}

func (s *AdminService) ListPayments(ctx context.Context, page, pageSize int) ([]*model.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.paymentRepo.ListAll(ctx, page, pageSize)
}

func (s *AdminService) FinancialSummary(ctx context.Context) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	stats, err := s.adminRepo.DashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get dashboard stats: %w", err)
	}

	summary["total_revenue"] = stats["total_revenue"]
	summary["active_contracts"] = stats["active_contracts"]
	summary["pending_disputes"] = stats["pending_disputes"]

	return summary, nil
}
