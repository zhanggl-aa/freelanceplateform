package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type ReviewRepository struct {
	db *DB
}

func NewReviewRepository(db *DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, projectID, contractID, reviewerID, revieweeID string, qualityRating, commRating, timelinessRating int, overallRating float64, comment *string, isPublic bool) (*model.Review, error) {
	var rv model.Review
	query := `
		INSERT INTO reviews (project_id, contract_id, reviewer_id, reviewee_id, quality_rating, communication_rating, timeliness_rating, overall_rating, comment, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, project_id, contract_id, reviewer_id, reviewee_id, quality_rating, communication_rating,
			timeliness_rating, overall_rating, comment, is_public, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		projectID, contractID, reviewerID, revieweeID,
		qualityRating, commRating, timelinessRating, overallRating, comment, isPublic,
	).Scan(
		&rv.ID, &rv.ProjectID, &rv.ContractID, &rv.ReviewerID, &rv.RevieweeID,
		&rv.QualityRating, &rv.CommunicationRating, &rv.TimelinessRating,
		&rv.OverallRating, &rv.Comment, &rv.IsPublic, &rv.CreatedAt, &rv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}
	return &rv, nil
}

func (r *ReviewRepository) GetByID(ctx context.Context, id string) (*model.Review, error) {
	var rv model.Review
	query := `
		SELECT r.id, r.project_id, r.contract_id, r.reviewer_id, r.reviewee_id,
			r.quality_rating, r.communication_rating, r.timeliness_rating, r.overall_rating,
			r.comment, r.is_public, r.created_at, r.updated_at,
			u.nickname AS reviewer_name, u.avatar_url AS reviewer_avatar
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&rv.ID, &rv.ProjectID, &rv.ContractID, &rv.ReviewerID, &rv.RevieweeID,
		&rv.QualityRating, &rv.CommunicationRating, &rv.TimelinessRating,
		&rv.OverallRating, &rv.Comment, &rv.IsPublic, &rv.CreatedAt, &rv.UpdatedAt,
		&rv.ReviewerName, &rv.ReviewerAvatar,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get review by id: %w", err)
	}
	return &rv, nil
}

func (r *ReviewRepository) GetByProject(ctx context.Context, projectID string) ([]*model.Review, error) {
	query := `
		SELECT r.id, r.project_id, r.contract_id, r.reviewer_id, r.reviewee_id,
			r.quality_rating, r.communication_rating, r.timeliness_rating, r.overall_rating,
			r.comment, r.is_public, r.created_at, r.updated_at,
			u.nickname AS reviewer_name, u.avatar_url AS reviewer_avatar
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.project_id = $1
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get reviews by project: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(
			&rv.ID, &rv.ProjectID, &rv.ContractID, &rv.ReviewerID, &rv.RevieweeID,
			&rv.QualityRating, &rv.CommunicationRating, &rv.TimelinessRating,
			&rv.OverallRating, &rv.Comment, &rv.IsPublic, &rv.CreatedAt, &rv.UpdatedAt,
			&rv.ReviewerName, &rv.ReviewerAvatar,
		); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		reviews = append(reviews, &rv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews: %w", err)
	}
	return reviews, nil
}

func (r *ReviewRepository) GetByReviewee(ctx context.Context, revieweeID string, page, pageSize int) ([]*model.Review, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM reviews WHERE reviewee_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, revieweeID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews by reviewee: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.project_id, r.contract_id, r.reviewer_id, r.reviewee_id,
			r.quality_rating, r.communication_rating, r.timeliness_rating, r.overall_rating,
			r.comment, r.is_public, r.created_at, r.updated_at,
			u.nickname AS reviewer_name, u.avatar_url AS reviewer_avatar
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.reviewee_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, revieweeID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get reviews by reviewee: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(
			&rv.ID, &rv.ProjectID, &rv.ContractID, &rv.ReviewerID, &rv.RevieweeID,
			&rv.QualityRating, &rv.CommunicationRating, &rv.TimelinessRating,
			&rv.OverallRating, &rv.Comment, &rv.IsPublic, &rv.CreatedAt, &rv.UpdatedAt,
			&rv.ReviewerName, &rv.ReviewerAvatar,
		); err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}
		reviews = append(reviews, &rv)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reviews: %w", err)
	}
	return reviews, total, nil
}

func (r *ReviewRepository) GetByReviewer(ctx context.Context, reviewerID string, page, pageSize int) ([]*model.Review, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM reviews WHERE reviewer_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, reviewerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews by reviewer: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.project_id, r.contract_id, r.reviewer_id, r.reviewee_id,
			r.quality_rating, r.communication_rating, r.timeliness_rating, r.overall_rating,
			r.comment, r.is_public, r.created_at, r.updated_at,
			u.nickname AS reviewer_name, u.avatar_url AS reviewer_avatar
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.reviewer_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, reviewerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get reviews by reviewer: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(
			&rv.ID, &rv.ProjectID, &rv.ContractID, &rv.ReviewerID, &rv.RevieweeID,
			&rv.QualityRating, &rv.CommunicationRating, &rv.TimelinessRating,
			&rv.OverallRating, &rv.Comment, &rv.IsPublic, &rv.CreatedAt, &rv.UpdatedAt,
			&rv.ReviewerName, &rv.ReviewerAvatar,
		); err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}
		reviews = append(reviews, &rv)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reviews: %w", err)
	}
	return reviews, total, nil
}

func (r *ReviewRepository) Update(ctx context.Context, id string, qualityRating, commRating, timelinessRating int, overallRating float64, comment *string) error {
	query := `
		UPDATE reviews
		SET quality_rating = $2, communication_rating = $3, timeliness_rating = $4,
			overall_rating = $5, comment = $6
		WHERE id = $1
	`
	tag, err := r.db.Pool.Exec(ctx, query, id, qualityRating, commRating, timelinessRating, overallRating, comment)
	if err != nil {
		return fmt.Errorf("update review: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reviews WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete review: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (r *ReviewRepository) GetAverageRating(ctx context.Context, revieweeID string) (float64, int, error) {
	var avg float64
	var count int
	query := `
		SELECT COALESCE(AVG(overall_rating), 0), COUNT(*)
		FROM reviews
		WHERE reviewee_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, revieweeID).Scan(&avg, &count)
	if err != nil {
		return 0, 0, fmt.Errorf("get average rating: %w", err)
	}
	return avg, count, nil
}
