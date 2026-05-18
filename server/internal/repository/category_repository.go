package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	db *DB
}

func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetTree(ctx context.Context) ([]*model.ProjectCategory, error) {
	query := `
		SELECT id::text, name, slug, description, icon_url, parent_id::text, sort_order, is_active, created_at
		FROM project_categories
		WHERE is_active = true
		ORDER BY sort_order ASC, name ASC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	all := make([]*model.ProjectCategory, 0)
	for rows.Next() {
		c := &model.ProjectCategory{}
		var parentID *string
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.IconURL, &parentID, &c.SortOrder, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		c.ParentID = parentID
		all = append(all, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	// Build tree: map of id -> category, attach children to parents
	idMap := make(map[string]*model.ProjectCategory, len(all))
	for _, c := range all {
		c.Children = []model.ProjectCategory{}
		idMap[c.ID] = c
	}

	roots := make([]*model.ProjectCategory, 0)
	for _, c := range all {
		if c.ParentID != nil && *c.ParentID != "" {
			if parent, ok := idMap[*c.ParentID]; ok {
				parent.Children = append(parent.Children, *c)
			} else {
				roots = append(roots, c)
			}
		} else {
			roots = append(roots, c)
		}
	}

	return roots, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*model.ProjectCategory, error) {
	query := `
		SELECT id::text, name, slug, description, icon_url, parent_id::text, sort_order, is_active, created_at
		FROM project_categories
		WHERE id = $1
	`
	c := &model.ProjectCategory{}
	var parentID *string
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.IconURL, &parentID, &c.SortOrder, &c.IsActive, &c.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	c.ParentID = parentID
	return c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, name, slug string, description, iconURL, parentID *string, sortOrder int) (*model.ProjectCategory, error) {
	query := `
		INSERT INTO project_categories (name, slug, description, icon_url, parent_id, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, name, slug, description, icon_url, parent_id::text, sort_order, is_active, created_at
	`
	c := &model.ProjectCategory{}
	var returnedParentID *string
	err := r.db.Pool.QueryRow(ctx, query, name, slug, description, iconURL, parentID, sortOrder).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.IconURL, &returnedParentID, &c.SortOrder, &c.IsActive, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	c.ParentID = returnedParentID
	return c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"name": true, "slug": true, "description": true, "icon_url": true,
		"parent_id": true, "sort_order": true, "is_active": true,
	}

	setClauses := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+1)
	argIdx := 1

	for col, val := range fields {
		if !allowed[col] {
			return fmt.Errorf("invalid column: %s", col)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE project_categories SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found: %s", id)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE project_categories SET is_active = false WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found: %s", id)
	}
	return nil
}
