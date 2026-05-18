package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type BookmarkService struct {
	bookmarkRepo *repository.BookmarkRepository
}

func NewBookmarkService(bookmarkRepo *repository.BookmarkRepository) *BookmarkService {
	return &BookmarkService{bookmarkRepo: bookmarkRepo}
}

func (s *BookmarkService) Create(ctx context.Context, userID, projectID string) (*model.Bookmark, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if projectID == "" {
		return nil, errors.New("project_id is required")
	}

	exists, err := s.bookmarkRepo.IsBookmarked(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("check existing bookmark: %w", err)
	}
	if exists {
		return nil, errors.New("project is already bookmarked")
	}

	created, err := s.bookmarkRepo.Create(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("create bookmark: %w", err)
	}

	return created, nil
}

func (s *BookmarkService) Delete(ctx context.Context, userID, projectID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}
	if projectID == "" {
		return errors.New("project_id is required")
	}

	return s.bookmarkRepo.Delete(ctx, userID, projectID)
}

func (s *BookmarkService) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Bookmark, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.bookmarkRepo.ListByUser(ctx, userID, page, pageSize)
}

func (s *BookmarkService) IsBookmarked(ctx context.Context, userID, projectID string) (bool, error) {
	if userID == "" || projectID == "" {
		return false, nil
	}

	return s.bookmarkRepo.IsBookmarked(ctx, userID, projectID)
}
