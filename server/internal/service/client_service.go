package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type ClientService struct {
	clientRepo *repository.ClientRepository
}

func NewClientService(clientRepo *repository.ClientRepository) *ClientService {
	return &ClientService{clientRepo: clientRepo}
}

// CreateProfile creates a new client profile for the given user.
func (s *ClientService) CreateProfile(ctx context.Context, profile *model.ClientProfile) (*model.ClientProfile, error) {
	if profile.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	existing, err := s.clientRepo.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("check existing profile: %w", err)
	}
	if existing != nil {
		return nil, errors.New("client profile already exists for this user")
	}

	req := map[string]interface{}{}
	if profile.CompanyName != nil {
		req["company_name"] = *profile.CompanyName
	}
	if profile.CompanyLogoURL != nil {
		req["company_logo_url"] = *profile.CompanyLogoURL
	}
	if profile.CompanyWebsite != nil {
		req["company_website"] = *profile.CompanyWebsite
	}
	if profile.Industry != nil {
		req["industry"] = *profile.Industry
	}
	if profile.CompanySize != nil {
		req["company_size"] = *profile.CompanySize
	}

	created, err := s.clientRepo.Create(ctx, profile.UserID, req)
	if err != nil {
		return nil, fmt.Errorf("create client profile: %w", err)
	}

	return created, nil
}

// GetProfile returns a client profile by its primary key.
func (s *ClientService) GetProfile(ctx context.Context, id string) (*model.ClientProfile, error) {
	profile, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get client profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("client profile not found")
	}
	return profile, nil
}

// GetByUserID returns a client profile by the owning user's ID.
// Returns (nil, nil) when no profile exists for the given user.
func (s *ClientService) GetByUserID(ctx context.Context, userID string) (*model.ClientProfile, error) {
	profile, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get client profile by user id: %w", err)
	}
	return profile, nil
}

// UpdateProfile updates the editable fields of a client profile.
func (s *ClientService) UpdateProfile(ctx context.Context, profile *model.ClientProfile) (*model.ClientProfile, error) {
	if profile.ID == "" {
		return nil, errors.New("profile id is required")
	}

	existing, err := s.clientRepo.GetByID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("check profile existence: %w", err)
	}
	if existing == nil {
		return nil, errors.New("client profile not found")
	}

	fields := map[string]interface{}{}
	if profile.CompanyName != nil {
		fields["company_name"] = *profile.CompanyName
	}
	if profile.CompanyLogoURL != nil {
		fields["company_logo_url"] = *profile.CompanyLogoURL
	}
	if profile.CompanyWebsite != nil {
		fields["company_website"] = *profile.CompanyWebsite
	}
	if profile.Industry != nil {
		fields["industry"] = *profile.Industry
	}
	if profile.CompanySize != nil {
		fields["company_size"] = *profile.CompanySize
	}

	updated, err := s.clientRepo.Update(ctx, profile.ID, fields)
	if err != nil {
		return nil, fmt.Errorf("update client profile: %w", err)
	}

	return updated, nil
}

// Verify marks a client profile as verified (admin action).
func (s *ClientService) Verify(ctx context.Context, clientID string) error {
	if clientID == "" {
		return errors.New("client_id is required")
	}

	profile, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("get client profile: %w", err)
	}
	if profile == nil {
		return errors.New("client profile not found")
	}
	if profile.Verified {
		return errors.New("client is already verified")
	}

	err = s.clientRepo.Verify(ctx, clientID)
	if err != nil {
		return fmt.Errorf("verify client: %w", err)
	}

	return nil
}
