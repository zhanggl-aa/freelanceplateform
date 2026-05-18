package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type DisputeRepository struct {
	db *DB
}

func NewDisputeRepository(db *DB) *DisputeRepository {
	return &DisputeRepository{db: db}
}

func (r *DisputeRepository) Create(ctx context.Context, contractID, milestoneID, reporterID, reportedID, reason string, evidenceURLs []string) (*model.Dispute, error) {
	var d model.Dispute
	var milestoneIDPtr *string
	if milestoneID != "" {
		milestoneIDPtr = &milestoneID
	}

	evidenceJSON, err := json.Marshal(evidenceURLs)
	if err != nil {
		return nil, fmt.Errorf("marshal evidence urls: %w", err)
	}

	query := `
		INSERT INTO disputes (contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls,
			status, resolution, resolution_type, resolved_by, resolved_at, created_at, updated_at
	`
	var evidenceRaw string
	err = r.db.Pool.QueryRow(ctx, query,
		contractID, milestoneIDPtr, reporterID, reportedID, reason, string(evidenceJSON),
	).Scan(
		&d.ID, &d.ContractID, &d.MilestoneID, &d.ReporterID, &d.ReportedID, &d.Reason, &evidenceRaw,
		&d.Status, &d.Resolution, &d.ResolutionType, &d.ResolvedBy, &d.ResolvedAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create dispute: %w", err)
	}

	if err := json.Unmarshal([]byte(evidenceRaw), &d.EvidenceURLs); err != nil {
		return nil, fmt.Errorf("unmarshal evidence urls: %w", err)
	}

	return &d, nil
}

func (r *DisputeRepository) GetByID(ctx context.Context, id string) (*model.Dispute, error) {
	var d model.Dispute
	var evidenceRaw string
	query := `
		SELECT id, contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls,
			status, resolution, resolution_type, resolved_by, resolved_at, created_at, updated_at
		FROM disputes
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.ContractID, &d.MilestoneID, &d.ReporterID, &d.ReportedID, &d.Reason, &evidenceRaw,
		&d.Status, &d.Resolution, &d.ResolutionType, &d.ResolvedBy, &d.ResolvedAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get dispute by id: %w", err)
	}

	if err := json.Unmarshal([]byte(evidenceRaw), &d.EvidenceURLs); err != nil {
		return nil, fmt.Errorf("unmarshal evidence urls: %w", err)
	}

	return &d, nil
}

func (r *DisputeRepository) ListByContract(ctx context.Context, contractID string) ([]*model.Dispute, error) {
	query := `
		SELECT id, contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls,
			status, resolution, resolution_type, resolved_by, resolved_at, created_at, updated_at
		FROM disputes
		WHERE contract_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, contractID)
	if err != nil {
		return nil, fmt.Errorf("list disputes by contract: %w", err)
	}
	defer rows.Close()

	var disputes []*model.Dispute
	for rows.Next() {
		var d model.Dispute
		var evidenceRaw string
		if err := rows.Scan(
			&d.ID, &d.ContractID, &d.MilestoneID, &d.ReporterID, &d.ReportedID, &d.Reason, &evidenceRaw,
			&d.Status, &d.Resolution, &d.ResolutionType, &d.ResolvedBy, &d.ResolvedAt, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dispute: %w", err)
		}
		if err := json.Unmarshal([]byte(evidenceRaw), &d.EvidenceURLs); err != nil {
			return nil, fmt.Errorf("unmarshal evidence urls: %w", err)
		}
		disputes = append(disputes, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate disputes: %w", err)
	}
	return disputes, nil
}

func (r *DisputeRepository) List(ctx context.Context, status string, page, pageSize int) ([]*model.Dispute, int64, error) {
	var total int64
	var countArgs []interface{}
	countQuery := `SELECT COUNT(*) FROM disputes`
	if status != "" {
		countQuery += ` WHERE status = $1`
		countArgs = append(countArgs, status)
	}
	if err := r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count disputes: %w", err)
	}

	offset := (page - 1) * pageSize
	var query string
	var args []interface{}
	if status != "" {
		query = `
			SELECT id, contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls,
				status, resolution, resolution_type, resolved_by, resolved_at, created_at, updated_at
			FROM disputes
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{status, pageSize, offset}
	} else {
		query = `
			SELECT id, contract_id, milestone_id, reporter_id, reported_id, reason, evidence_urls,
				status, resolution, resolution_type, resolved_by, resolved_at, created_at, updated_at
			FROM disputes
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{pageSize, offset}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list disputes: %w", err)
	}
	defer rows.Close()

	var disputes []*model.Dispute
	for rows.Next() {
		var d model.Dispute
		var evidenceRaw string
		if err := rows.Scan(
			&d.ID, &d.ContractID, &d.MilestoneID, &d.ReporterID, &d.ReportedID, &d.Reason, &evidenceRaw,
			&d.Status, &d.Resolution, &d.ResolutionType, &d.ResolvedBy, &d.ResolvedAt, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan dispute: %w", err)
		}
		if err := json.Unmarshal([]byte(evidenceRaw), &d.EvidenceURLs); err != nil {
			return nil, 0, fmt.Errorf("unmarshal evidence urls: %w", err)
		}
		disputes = append(disputes, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate disputes: %w", err)
	}
	return disputes, total, nil
}

func (r *DisputeRepository) UpdateResolution(ctx context.Context, id string, resolution string, resolutionType string, resolvedBy string) error {
	query := `
		UPDATE disputes
		SET resolution = $2, resolution_type = $3, resolved_by = $4, resolved_at = now(), status = 'resolved'
		WHERE id = $1
	`
	tag, err := r.db.Pool.Exec(ctx, query, id, resolution, resolutionType, resolvedBy)
	if err != nil {
		return fmt.Errorf("update dispute resolution: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("dispute not found")
	}
	return nil
}

func (r *DisputeRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE disputes SET status = $2 WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("update dispute status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("dispute not found")
	}
	return nil
}
