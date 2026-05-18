package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type DeveloperService struct {
	developerRepo *repository.DeveloperRepository
}

func NewDeveloperService(developerRepo *repository.DeveloperRepository) *DeveloperService {
	return &DeveloperService{developerRepo: developerRepo}
}

// CreateProfile creates a new developer profile for the given user.
func (s *DeveloperService) CreateProfile(ctx context.Context, userID string, fields map[string]interface{}) (*model.DeveloperProfile, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	existing, err := s.developerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check existing profile: %w", err)
	}
	if existing != nil {
		return nil, errors.New("developer profile already exists for this user")
	}

	if _, ok := fields["availability"]; !ok {
		fields["availability"] = "available"
	}

	created, err := s.developerRepo.Create(ctx, userID, fields)
	if err != nil {
		return nil, fmt.Errorf("create developer profile: %w", err)
	}

	return created, nil
}

// GetProfile returns a developer profile by its primary key.
func (s *DeveloperService) GetProfile(ctx context.Context, id string) (*model.DeveloperProfile, error) {
	profile, err := s.developerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get developer profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("developer profile not found")
	}
	return profile, nil
}

// GetByUserID returns a developer profile by the owning user's ID.
// Returns (nil, nil) when no profile exists for the given user.
func (s *DeveloperService) GetByUserID(ctx context.Context, userID string) (*model.DeveloperProfile, error) {
	profile, err := s.developerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get developer profile by user id: %w", err)
	}
	return profile, nil
}

// UpdateProfile updates the editable fields of a developer profile.
func (s *DeveloperService) UpdateProfile(ctx context.Context, id string, fields map[string]interface{}) (*model.DeveloperProfile, error) {
	if id == "" {
		return nil, errors.New("profile id is required")
	}
	if len(fields) == 0 {
		return nil, errors.New("no fields to update")
	}

	existing, err := s.developerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("check profile existence: %w", err)
	}
	if existing == nil {
		return nil, errors.New("developer profile not found")
	}

	updated, err := s.developerRepo.Update(ctx, id, fields)
	if err != nil {
		return nil, fmt.Errorf("update developer profile: %w", err)
	}

	return updated, nil
}

// AddSkill adds a new skill to a developer profile.
func (s *DeveloperService) AddSkill(ctx context.Context, developerID, skillName, proficiency string, yearsExperience *int) (*model.DeveloperSkill, error) {
	if developerID == "" {
		return nil, errors.New("developer_id is required")
	}
	if skillName == "" {
		return nil, errors.New("skill_name is required")
	}
	if proficiency == "" {
		proficiency = "intermediate"
	}

	validProficiencies := map[string]bool{"beginner": true, "intermediate": true, "advanced": true, "expert": true}
	if !validProficiencies[proficiency] {
		return nil, errors.New("proficiency must be one of: beginner, intermediate, advanced, expert")
	}

	created, err := s.developerRepo.AddSkill(ctx, developerID, skillName, proficiency, yearsExperience)
	if err != nil {
		return nil, fmt.Errorf("add skill: %w", err)
	}

	return created, nil
}

// UpdateSkill modifies an existing skill.
func (s *DeveloperService) UpdateSkill(ctx context.Context, skillID, proficiency string, yearsExperience *int) error {
	if skillID == "" {
		return errors.New("skill id is required")
	}

	if proficiency != "" {
		validProficiencies := map[string]bool{"beginner": true, "intermediate": true, "advanced": true, "expert": true}
		if !validProficiencies[proficiency] {
			return errors.New("proficiency must be one of: beginner, intermediate, advanced, expert")
		}
	}

	return s.developerRepo.UpdateSkill(ctx, skillID, proficiency, yearsExperience)
}

// DeleteSkill removes a skill from a developer profile.
func (s *DeveloperService) DeleteSkill(ctx context.Context, skillID string) error {
	if skillID == "" {
		return errors.New("skill id is required")
	}

	return s.developerRepo.DeleteSkill(ctx, skillID)
}

// ListSkills returns all skills for a developer.
func (s *DeveloperService) ListSkills(ctx context.Context, developerID string) ([]*model.DeveloperSkill, error) {
	if developerID == "" {
		return nil, errors.New("developer_id is required")
	}

	skills, err := s.developerRepo.ListSkills(ctx, developerID)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}

	return skills, nil
}

// AddPortfolio adds a portfolio item to a developer profile.
func (s *DeveloperService) AddPortfolio(ctx context.Context, developerID string, portfolio *model.DeveloperPortfolio) (*model.DeveloperPortfolio, error) {
	if developerID == "" {
		return nil, errors.New("developer_id is required")
	}
	if portfolio.Title == "" {
		return nil, errors.New("title is required")
	}

	if portfolio.TechStack == nil {
		portfolio.TechStack = []string{}
	}
	if portfolio.ImageURLs == nil {
		portfolio.ImageURLs = []string{}
	}

	created, err := s.developerRepo.AddPortfolio(ctx, developerID, portfolio)
	if err != nil {
		return nil, fmt.Errorf("add portfolio: %w", err)
	}

	return created, nil
}

// UpdatePortfolio modifies an existing portfolio item.
func (s *DeveloperService) UpdatePortfolio(ctx context.Context, portfolioID string, portfolio *model.DeveloperPortfolio) (*model.DeveloperPortfolio, error) {
	if portfolioID == "" {
		return nil, errors.New("portfolio id is required")
	}

	updated, err := s.developerRepo.UpdatePortfolio(ctx, portfolioID, portfolio)
	if err != nil {
		return nil, fmt.Errorf("update portfolio: %w", err)
	}

	return updated, nil
}

// DeletePortfolio removes a portfolio item.
func (s *DeveloperService) DeletePortfolio(ctx context.Context, portfolioID string) error {
	if portfolioID == "" {
		return errors.New("portfolio id is required")
	}

	return s.developerRepo.DeletePortfolio(ctx, portfolioID)
}

// ListPortfolio returns all portfolio items for a developer.
func (s *DeveloperService) ListPortfolio(ctx context.Context, developerID string) ([]*model.DeveloperPortfolio, error) {
	if developerID == "" {
		return nil, errors.New("developer_id is required")
	}

	items, err := s.developerRepo.ListPortfolio(ctx, developerID)
	if err != nil {
		return nil, fmt.Errorf("list portfolio: %w", err)
	}

	return items, nil
}

// Search finds developers matching a keyword, skills, or availability criteria.
func (s *DeveloperService) Search(ctx context.Context, keyword string, skills []string, availability string, minRate, maxRate *float64, page, pageSize int) ([]*model.DeveloperProfile, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	skillFilter := ""
	if len(skills) > 0 {
		skillFilter = strings.Join(skills, ",")
	}
	if keyword != "" && skillFilter == "" {
		skillFilter = keyword
	}

	results, total, err := s.developerRepo.Search(ctx, skillFilter, minRate, maxRate, availability, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("search developers: %w", err)
	}

	return results, total, nil
}

// Verify marks a developer profile as verified (admin action).
func (s *DeveloperService) Verify(ctx context.Context, developerID string) error {
	if developerID == "" {
		return errors.New("developer_id is required")
	}

	profile, err := s.developerRepo.GetByID(ctx, developerID)
	if err != nil {
		return fmt.Errorf("get developer profile: %w", err)
	}
	if profile == nil {
		return errors.New("developer profile not found")
	}
	if profile.Verified {
		return errors.New("developer is already verified")
	}

	now := time.Now().Format(time.RFC3339)
	_, err = s.developerRepo.Update(ctx, developerID, map[string]interface{}{
		"verified":   true,
		"verified_at": now,
	})
	if err != nil {
		return fmt.Errorf("verify developer: %w", err)
	}

	return nil
}
