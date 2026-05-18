package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type FileRepository struct {
	db *DB
}

func NewFileRepository(db *DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, userID, fileName, filePath string, fileSize int64, fileType, mimeType, storageType string, entityType, entityID *string) (*model.FileAttachment, error) {
	var f model.FileAttachment
	query := `
		INSERT INTO file_attachments (user_id, file_name, file_path, file_size, file_type, mime_type, storage_type, entity_type, entity_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, file_name, file_path, file_size, file_type, mime_type, storage_type, entity_type, entity_id, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		userID, fileName, filePath, fileSize, fileType, mimeType, storageType, entityType, entityID,
	).Scan(
		&f.ID, &f.UserID, &f.FileName, &f.FilePath, &f.FileSize, &f.FileType,
		&f.MimeType, &f.StorageType, &f.EntityType, &f.EntityID, &f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create file attachment: %w", err)
	}
	return &f, nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*model.FileAttachment, error) {
	var f model.FileAttachment
	query := `
		SELECT id, user_id, file_name, file_path, file_size, file_type, mime_type, storage_type, entity_type, entity_id, created_at
		FROM file_attachments
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.UserID, &f.FileName, &f.FilePath, &f.FileSize, &f.FileType,
		&f.MimeType, &f.StorageType, &f.EntityType, &f.EntityID, &f.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get file attachment by id: %w", err)
	}
	return &f, nil
}

func (r *FileRepository) ListByEntity(ctx context.Context, entityType, entityID string) ([]*model.FileAttachment, error) {
	query := `
		SELECT id, user_id, file_name, file_path, file_size, file_type, mime_type, storage_type, entity_type, entity_id, created_at
		FROM file_attachments
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("list file attachments by entity: %w", err)
	}
	defer rows.Close()

	var files []*model.FileAttachment
	for rows.Next() {
		var f model.FileAttachment
		if err := rows.Scan(
			&f.ID, &f.UserID, &f.FileName, &f.FilePath, &f.FileSize, &f.FileType,
			&f.MimeType, &f.StorageType, &f.EntityType, &f.EntityID, &f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan file attachment: %w", err)
		}
		files = append(files, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate file attachments: %w", err)
	}
	return files, nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM file_attachments WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete file attachment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("file attachment not found")
	}
	return nil
}
