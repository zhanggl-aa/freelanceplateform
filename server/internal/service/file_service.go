package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

// StorageProvider defines the interface for file storage backends.
type StorageProvider interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
	GetURL(key string) string
	Delete(ctx context.Context, key string) error
}

// LocalStorageProvider implements StorageProvider for local filesystem storage.
type LocalStorageProvider struct {
	basePath string
	baseURL  string
}

func NewLocalStorageProvider(basePath, baseURL string) *LocalStorageProvider {
	return &LocalStorageProvider{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Upload saves a file to the local filesystem under basePath/key.
func (p *LocalStorageProvider) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	if key == "" {
		return "", errors.New("storage key is required")
	}

	fullPath := filepath.Join(p.basePath, key)
	dir := filepath.Dir(fullPath)

	// Create directory if it does not exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return p.GetURL(key), nil
}

// GetURL returns the URL for a stored file.
func (p *LocalStorageProvider) GetURL(key string) string {
	if p.baseURL != "" {
		return p.baseURL + "/" + key
	}
	return "/" + key
}

// Delete removes a file from the local filesystem.
func (p *LocalStorageProvider) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("storage key is required")
	}

	fullPath := filepath.Join(p.basePath, key)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

// FileService handles file upload, retrieval, and deletion.
type FileService struct {
	fileRepo *repository.FileRepository
	storage  StorageProvider
}

func NewFileService(fileRepo *repository.FileRepository, storage StorageProvider) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

// Upload saves a file to storage and creates a DB record.
// The key is generated as: {userID}/{entityType}/{entityID}/{filename} or {userID}/{filename} if no entity.
func (s *FileService) Upload(ctx context.Context, userID, filename string, reader io.Reader, contentType string, fileSize int64, entityType, entityID *string) (*model.FileAttachment, string, error) {
	if userID == "" {
		return nil, "", errors.New("user_id is required")
	}
	if filename == "" {
		return nil, "", errors.New("filename is required")
	}
	if reader == nil {
		return nil, "", errors.New("file data is required")
	}

	// Sanitize filename
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, " ", "_")

	// Generate storage key
	var key string
	if entityType != nil && *entityType != "" && entityID != nil && *entityID != "" {
		key = fmt.Sprintf("%s/%s/%s/%s", userID, *entityType, *entityID, filename)
	} else {
		key = fmt.Sprintf("%s/%s", userID, filename)
	}

	// Detect file type from extension
	fileType := detectFileType(filename)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to storage
	url, err := s.storage.Upload(ctx, key, reader, contentType)
	if err != nil {
		return nil, "", fmt.Errorf("upload to storage: %w", err)
	}

	// Create DB record
	created, err := s.fileRepo.Create(ctx, userID, filename, key, fileSize, fileType, contentType, "local", entityType, entityID)
	if err != nil {
		// Attempt to clean up the uploaded file if DB record fails
		_ = s.storage.Delete(ctx, key)
		return nil, "", fmt.Errorf("create file record: %w", err)
	}

	return created, url, nil
}

// GetByID returns a file attachment record by its ID.
func (s *FileService) GetByID(ctx context.Context, id string) (*model.FileAttachment, error) {
	if id == "" {
		return nil, errors.New("file id is required")
	}

	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}
	return file, nil
}

// GetURL returns the URL for a file by its ID.
func (s *FileService) GetURL(ctx context.Context, id string) (string, error) {
	file, err := s.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	return s.storage.GetURL(file.FilePath), nil
}

// Delete removes a file from storage and deletes its DB record.
func (s *FileService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("file id is required")
	}

	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}
	if file == nil {
		return errors.New("file not found")
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, file.FilePath); err != nil {
		return fmt.Errorf("delete from storage: %w", err)
	}

	// Delete DB record
	return s.fileRepo.Delete(ctx, id)
}

// detectFileType returns a broad file type category based on extension.
func detectFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp":
		return "image"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "document"
	case ".xls", ".xlsx":
		return "spreadsheet"
	case ".ppt", ".pptx":
		return "presentation"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "archive"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "video"
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "audio"
	case ".txt", ".md", ".csv":
		return "text"
	default:
		return "other"
	}
}
