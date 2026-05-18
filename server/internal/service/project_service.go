package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

// Create creates a new project in draft status.
func (s *ProjectService) Create(ctx context.Context, project *model.Project) (*model.Project, error) {
	if project.ClientID == "" {
		return nil, errors.New("client_id is required")
	}
	if project.CategoryID == "" {
		return nil, errors.New("category_id is required")
	}
	if project.Title == "" {
		return nil, errors.New("title is required")
	}
	if project.Description == "" {
		return nil, errors.New("description is required")
	}
	if project.BudgetType == "" {
		project.BudgetType = "fixed"
	}

	validBudgetTypes := map[string]bool{"fixed": true, "hourly": true}
	if !validBudgetTypes[project.BudgetType] {
		return nil, errors.New("budget_type must be 'fixed' or 'hourly'")
	}

	if project.BudgetMin != nil && project.BudgetMax != nil && *project.BudgetMin > *project.BudgetMax {
		return nil, errors.New("budget_min cannot exceed budget_max")
	}

	if project.TechStack == nil {
		project.TechStack = []string{}
	}

	var deadline *string
	if project.Deadline != nil {
		d := project.Deadline.Format("2006-01-02")
		deadline = &d
	}

	created, err := s.projectRepo.Create(ctx, project.ClientID, project.CategoryID, project.Title, project.Description, project.BudgetMin, project.BudgetMax, project.BudgetType, deadline, project.TechStack)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return created, nil
}

// GetByID returns a single project by its ID.
func (s *ProjectService) GetByID(ctx context.Context, id string) (*model.Project, error) {
	if id == "" {
		return nil, errors.New("project id is required")
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return project, nil
}

// Update modifies an existing project. Only draft or open projects can be updated.
func (s *ProjectService) Update(ctx context.Context, project *model.Project) (*model.Project, error) {
	if project.ID == "" {
		return nil, errors.New("project id is required")
	}

	existing, err := s.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if existing == nil {
		return nil, errors.New("project not found")
	}

	if existing.Status != "draft" && existing.Status != "open" {
		return nil, errors.New("only draft or open projects can be updated")
	}

	if project.BudgetType != "" {
		validBudgetTypes := map[string]bool{"fixed": true, "hourly": true}
		if !validBudgetTypes[project.BudgetType] {
			return nil, errors.New("budget_type must be 'fixed' or 'hourly'")
		}
	}

	if project.BudgetMin != nil && project.BudgetMax != nil && *project.BudgetMin > *project.BudgetMax {
		return nil, errors.New("budget_min cannot exceed budget_max")
	}

	fields := map[string]interface{}{}
	if project.Title != "" {
		fields["title"] = project.Title
	}
	if project.Description != "" {
		fields["description"] = project.Description
	}
	if project.BudgetMin != nil {
		fields["budget_min"] = *project.BudgetMin
	}
	if project.BudgetMax != nil {
		fields["budget_max"] = *project.BudgetMax
	}
	if project.BudgetType != "" {
		fields["budget_type"] = project.BudgetType
	}
	if project.Deadline != nil {
		fields["deadline"] = project.Deadline
	}
	if project.TechStack != nil {
		fields["tech_stack"] = project.TechStack
	}
	if project.CategoryID != "" {
		fields["category_id"] = project.CategoryID
	}
	if project.BidDeadline != nil {
		fields["bid_deadline"] = project.BidDeadline
	}
	if project.AssignedDeveloperID != nil {
		fields["assigned_developer_id"] = *project.AssignedDeveloperID
	}
	if project.Featured != existing.Featured {
		fields["featured"] = project.Featured
	}

	err = s.projectRepo.Update(ctx, project.ID, fields)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	updated, err := s.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("get updated project: %w", err)
	}
	return updated, nil
}

// Delete removes a project. Only draft projects can be deleted.
func (s *ProjectService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("project id is required")
	}

	existing, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}
	if existing == nil {
		return errors.New("project not found")
	}

	if existing.Status != "draft" {
		return errors.New("only draft projects can be deleted")
	}

	return s.projectRepo.Delete(ctx, id)
}

// Publish changes a project from draft to open status.
func (s *ProjectService) Publish(ctx context.Context, id string) (*model.Project, error) {
	if id == "" {
		return nil, errors.New("project id is required")
	}

	existing, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if existing == nil {
		return nil, errors.New("project not found")
	}

	if existing.Status != "draft" {
		return nil, errors.New("only draft projects can be published")
	}

	err = s.projectRepo.Update(ctx, id, map[string]interface{}{"status": "open"})
	if err != nil {
		return nil, fmt.Errorf("publish project: %w", err)
	}

	updated, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated project: %w", err)
	}
	return updated, nil
}

// Close marks an open project as closed.
func (s *ProjectService) Close(ctx context.Context, id string) (*model.Project, error) {
	if id == "" {
		return nil, errors.New("project id is required")
	}

	existing, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if existing == nil {
		return nil, errors.New("project not found")
	}

	if existing.Status != "open" {
		return nil, errors.New("only open projects can be closed")
	}

	err = s.projectRepo.Update(ctx, id, map[string]interface{}{"status": "closed"})
	if err != nil {
		return nil, fmt.Errorf("close project: %w", err)
	}

	updated, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated project: %w", err)
	}
	return updated, nil
}

// Search finds projects matching keyword, category, skills, or status filters.
func (s *ProjectService) Search(ctx context.Context, keyword string, categoryID string, skills []string, status string, page, pageSize int) ([]*model.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	keyword = strings.TrimSpace(keyword)

	// Default to searching published projects if no status specified
	if status == "" {
		status = "published"
	}

	var categoryIDPtr *string
	if categoryID != "" {
		categoryIDPtr = &categoryID
	}
	var keywordPtr *string
	if keyword != "" {
		keywordPtr = &keyword
	}
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	results, total, err := s.projectRepo.Search(ctx, categoryIDPtr, statusPtr, keywordPtr, nil, nil, skills, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("search projects: %w", err)
	}

	return results, total, nil
}

// ListByClient returns all projects posted by a specific client.
func (s *ProjectService) ListByClient(ctx context.Context, clientID string, page, pageSize int) ([]*model.Project, int64, error) {
	if clientID == "" {
		return nil, 0, errors.New("client_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	results, total, err := s.projectRepo.ListByClient(ctx, clientID, nil, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects by client: %w", err)
	}

	return results, total, nil
}

// ListByDeveloper returns projects assigned to or completed by a specific developer.
func (s *ProjectService) ListByDeveloper(ctx context.Context, developerID string, page, pageSize int) ([]*model.Project, int64, error) {
	if developerID == "" {
		return nil, 0, errors.New("developer_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	results, total, err := s.projectRepo.ListByDeveloper(ctx, developerID, nil, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects by developer: %w", err)
	}

	return results, total, nil
}

// ListFeatured returns projects marked as featured.
func (s *ProjectService) ListFeatured(ctx context.Context, page, pageSize int) ([]*model.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	results, total, err := s.projectRepo.ListFeatured(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list featured projects: %w", err)
	}

	return results, total, nil
}

// IncrementView atomically increments the view count for a project.
func (s *ProjectService) IncrementView(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("project id is required")
	}

	return s.projectRepo.IncrementViewCount(ctx, id)
}
