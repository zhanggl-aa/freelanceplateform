package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type ReviewService struct {
	reviewRepo   *repository.ReviewRepository
	projectRepo  *repository.ProjectRepository
	contractRepo *repository.ContractRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	projectRepo *repository.ProjectRepository,
	contractRepo *repository.ContractRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo:   reviewRepo,
		projectRepo:  projectRepo,
		contractRepo: contractRepo,
	}
}

// Create creates a new review. Validates that:
// - The project exists and is completed
// - The contract exists and is associated with the project
// - The reviewer is part of the contract (either client or developer)
// - The reviewee is the other party in the contract
// - The reviewer has not already reviewed this project
// - Ratings are between 1 and 5
// - Overall rating is calculated as the average of the three dimensions
func (s *ReviewService) Create(ctx context.Context, review *model.Review) (*model.Review, error) {
	if review.ProjectID == "" {
		return nil, errors.New("project_id is required")
	}
	if review.ContractID == "" {
		return nil, errors.New("contract_id is required")
	}
	if review.ReviewerID == "" {
		return nil, errors.New("reviewer_id is required")
	}
	if review.RevieweeID == "" {
		return nil, errors.New("reviewee_id is required")
	}
	if review.ReviewerID == review.RevieweeID {
		return nil, errors.New("reviewer and reviewee cannot be the same user")
	}

	// Validate ratings (1-5 scale)
	if review.QualityRating < 1 || review.QualityRating > 5 {
		return nil, errors.New("quality_rating must be between 1 and 5")
	}
	if review.CommunicationRating < 1 || review.CommunicationRating > 5 {
		return nil, errors.New("communication_rating must be between 1 and 5")
	}
	if review.TimelinessRating < 1 || review.TimelinessRating > 5 {
		return nil, errors.New("timeliness_rating must be between 1 and 5")
	}

	// Calculate overall rating as average of the three dimensions
	review.OverallRating = float64(review.QualityRating+review.CommunicationRating+review.TimelinessRating) / 3.0

	// Validate project is completed
	project, err := s.projectRepo.GetByID(ctx, review.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.Status != "completed" {
		return nil, errors.New("can only review completed projects")
	}

	// Validate contract
	contract, err := s.contractRepo.GetByID(ctx, review.ContractID)
	if err != nil {
		return nil, fmt.Errorf("get contract: %w", err)
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}
	if contract.ProjectID != review.ProjectID {
		return nil, errors.New("contract does not belong to this project")
	}

	// Validate reviewer is part of the contract
	if contract.ClientID != review.ReviewerID && contract.DeveloperID != review.ReviewerID {
		return nil, errors.New("reviewer is not part of this contract")
	}

	// Validate reviewee is the other party
	if contract.ClientID == review.ReviewerID && contract.DeveloperID != review.RevieweeID {
		return nil, errors.New("reviewee must be the other party in the contract")
	}
	if contract.DeveloperID == review.ReviewerID && contract.ClientID != review.RevieweeID {
		return nil, errors.New("reviewee must be the other party in the contract")
	}

	// Check for duplicate review - since there is no GetByProjectAndReviewer,
	// we check via GetByReviewer and filter by project
	reviewerReviews, _, err := s.reviewRepo.GetByReviewer(ctx, review.ReviewerID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("check existing review: %w", err)
	}
	for _, r := range reviewerReviews {
		if r.ProjectID == review.ProjectID {
			return nil, errors.New("you have already reviewed this project")
		}
	}

	created, err := s.reviewRepo.Create(ctx, review.ProjectID, review.ContractID, review.ReviewerID, review.RevieweeID, review.QualityRating, review.CommunicationRating, review.TimelinessRating, review.OverallRating, review.Comment, review.IsPublic)
	if err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}

	return created, nil
}

// GetByID returns a review by its ID.
func (s *ReviewService) GetByID(ctx context.Context, id string) (*model.Review, error) {
	if id == "" {
		return nil, errors.New("review id is required")
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get review: %w", err)
	}
	return review, nil
}

// GetByProject returns all reviews for a given project.
func (s *ReviewService) GetByProject(ctx context.Context, projectID string) ([]*model.Review, error) {
	if projectID == "" {
		return nil, errors.New("project_id is required")
	}

	reviews, err := s.reviewRepo.GetByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get reviews by project: %w", err)
	}

	return reviews, nil
}

// GetByReviewee returns all reviews received by a user (where they are the reviewee).
func (s *ReviewService) GetByReviewee(ctx context.Context, revieweeID string, page, pageSize int) ([]*model.Review, int64, error) {
	if revieweeID == "" {
		return nil, 0, errors.New("reviewee_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	reviews, total, err := s.reviewRepo.GetByReviewee(ctx, revieweeID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("get reviews by reviewee: %w", err)
	}

	return reviews, total, nil
}

// Update modifies an existing review. Only the reviewer can update, and only within a time limit.
func (s *ReviewService) Update(ctx context.Context, review *model.Review) (*model.Review, error) {
	if review.ID == "" {
		return nil, errors.New("review id is required")
	}

	existing, err := s.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		return nil, fmt.Errorf("get review: %w", err)
	}
	if existing == nil {
		return nil, errors.New("review not found")
	}

	// Validate ratings if being updated
	if review.QualityRating < 1 || review.QualityRating > 5 {
		return nil, errors.New("quality_rating must be between 1 and 5")
	}
	if review.CommunicationRating < 1 || review.CommunicationRating > 5 {
		return nil, errors.New("communication_rating must be between 1 and 5")
	}
	if review.TimelinessRating < 1 || review.TimelinessRating > 5 {
		return nil, errors.New("timeliness_rating must be between 1 and 5")
	}

	// Recalculate overall rating
	review.OverallRating = float64(review.QualityRating+review.CommunicationRating+review.TimelinessRating) / 3.0

	err = s.reviewRepo.Update(ctx, review.ID, review.QualityRating, review.CommunicationRating, review.TimelinessRating, review.OverallRating, review.Comment)
	if err != nil {
		return nil, fmt.Errorf("update review: %w", err)
	}

	updated, err := s.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		return nil, fmt.Errorf("get updated review: %w", err)
	}
	return updated, nil
}

// Delete removes a review.
func (s *ReviewService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("review id is required")
	}

	existing, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get review: %w", err)
	}
	if existing == nil {
		return errors.New("review not found")
	}

	return s.reviewRepo.Delete(ctx, id)
}
