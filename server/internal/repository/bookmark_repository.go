package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
)

type BookmarkRepository struct {
	db *DB
}

func NewBookmarkRepository(db *DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Create(ctx context.Context, userID, projectID string) (*model.Bookmark, error) {
	var b model.Bookmark
	query := `
		INSERT INTO bookmarks (user_id, project_id)
		VALUES ($1, $2)
		RETURNING id, user_id, project_id, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, userID, projectID).Scan(
		&b.ID, &b.UserID, &b.ProjectID, &b.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create bookmark: %w", err)
	}
	return &b, nil
}

func (r *BookmarkRepository) Delete(ctx context.Context, userID, projectID string) error {
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND project_id = $2`
	tag, err := r.db.Pool.Exec(ctx, query, userID, projectID)
	if err != nil {
		return fmt.Errorf("delete bookmark: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bookmark not found")
	}
	return nil
}

func (r *BookmarkRepository) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Bookmark, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bookmarks by user: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT b.id, b.user_id, b.project_id, b.created_at
		FROM bookmarks b
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookmarks by user: %w", err)
	}
	defer rows.Close()

	var bookmarks []*model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.ProjectID, &b.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate bookmarks: %w", err)
	}
	return bookmarks, total, nil
}

func (r *BookmarkRepository) IsBookmarked(ctx context.Context, userID, projectID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND project_id = $2)`
	err := r.db.Pool.QueryRow(ctx, query, userID, projectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check bookmark exists: %w", err)
	}
	return exists, nil
}
