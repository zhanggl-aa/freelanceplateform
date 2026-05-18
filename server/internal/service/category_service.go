package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// GetTree returns the active category tree (parent -> children hierarchy).
func (s *CategoryService) GetTree(ctx context.Context) ([]*model.ProjectCategory, error) {
	tree, err := s.categoryRepo.GetTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("get category tree: %w", err)
	}
	return tree, nil
}

// GetByID returns a single category by ID.
func (s *CategoryService) GetByID(ctx context.Context, id string) (*model.ProjectCategory, error) {
	if id == "" {
		return nil, errors.New("category id is required")
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	return category, nil
}

// Create adds a new category. Requires admin role (enforced at handler layer).
func (s *CategoryService) Create(ctx context.Context, name, slug string, description, iconURL, parentID *string, sortOrder int) (*model.ProjectCategory, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}
	if slug == "" {
		return nil, errors.New("category slug is required")
	}

	// Validate slug format: lowercase alphanumeric and hyphens only
	slug = strings.ToLower(slug)
	if !isValidSlug(slug) {
		return nil, errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}

	// If parent is specified, verify it exists
	if parentID != nil && *parentID != "" {
		parent, err := s.categoryRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("check parent category: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
	}

	category, err := s.categoryRepo.Create(ctx, name, slug, description, iconURL, parentID, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return category, nil
}

// Update modifies an existing category. Requires admin role (enforced at handler layer).
func (s *CategoryService) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if id == "" {
		return errors.New("category id is required")
	}
	if len(fields) == 0 {
		return errors.New("no fields to update")
	}

	// Validate slug if provided
	if slugVal, ok := fields["slug"]; ok {
		slug, ok := slugVal.(string)
		if !ok {
			return errors.New("slug must be a string")
		}
		slug = strings.ToLower(slug)
		if !isValidSlug(slug) {
			return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
		}
		fields["slug"] = slug
	}

	// Validate parent_id if provided
	if parentIDVal, ok := fields["parent_id"]; ok {
		if parentID, ok := parentIDVal.(*string); ok && parentID != nil && *parentID != "" {
			// Prevent self-reference
			if *parentID == id {
				return errors.New("category cannot be its own parent")
			}
			parent, err := s.categoryRepo.GetByID(ctx, *parentID)
			if err != nil {
				return fmt.Errorf("check parent category: %w", err)
			}
			if parent == nil {
				return errors.New("parent category not found")
			}
		}
	}

	return s.categoryRepo.Update(ctx, id, fields)
}

// Delete soft-deletes a category by setting is_active = false. Requires admin role.
func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id is required")
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get category: %w", err)
	}
	if category == nil {
		return errors.New("category not found")
	}
	if !category.IsActive {
		return errors.New("category is already inactive")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// isValidSlug checks that a slug contains only lowercase alphanumeric characters and hyphens.
func isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}
	for _, c := range slug {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}
