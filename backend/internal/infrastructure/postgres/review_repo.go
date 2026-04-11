package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rostmebel/backend/internal/domain/review"
)

type ReviewRepo struct {
	pool *pgxpool.Pool
}

func NewReviewRepo(pool *pgxpool.Pool) *ReviewRepo {
	return &ReviewRepo{pool: pool}
}

func (r *ReviewRepo) Create(ctx context.Context, rev *review.Review) error {
	images, _ := json.Marshal(rev.Images)
	query := `INSERT INTO reviews (project_id, order_id, rating, comment, images, status) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	
	return r.pool.QueryRow(ctx, query, 
		rev.ProjectID, rev.OrderID, rev.Rating, rev.Comment, images, rev.Status,
	).Scan(&rev.ID, &rev.CreatedAt, &rev.UpdatedAt)
}

func (r *ReviewRepo) GetByID(ctx context.Context, id int64) (*review.Review, error) {
	query := `
		SELECT r.id, r.project_id, r.order_id, r.rating, r.comment, r.images, r.status, r.created_at, r.updated_at, o.client_name, p.name
		FROM reviews r
		JOIN orders o ON r.order_id = o.id
		LEFT JOIN projects p ON r.project_id = p.id
		WHERE r.id = $1`
	
	var rev review.Review
	var images []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rev.ID, &rev.ProjectID, &rev.OrderID, &rev.Rating, &rev.Comment, &images, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt, &rev.ClientName, &rev.ProjectName,
	)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	json.Unmarshal(images, &rev.Images)
	return &rev, nil
}

func (r *ReviewRepo) List(ctx context.Context, f review.ListFilter) ([]*review.Review, int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("r.status = $%d", argCount))
		args = append(args, f.Status)
		argCount++
	}
	if f.ProjectID != nil {
		conditions = append(conditions, fmt.Sprintf("r.project_id = $%d", argCount))
		args = append(args, *f.ProjectID)
		argCount++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM reviews r %s", where)
	var total int64
	r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)

	query := fmt.Sprintf(`
		SELECT r.id, r.project_id, r.order_id, r.rating, r.comment, r.images, r.status, r.created_at, r.updated_at, o.client_name, p.name
		FROM reviews r
		JOIN orders o ON r.order_id = o.id
		LEFT JOIN projects p ON r.project_id = p.id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d`, where, argCount, argCount+1)
	
	args = append(args, f.Limit, f.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()

	var reviews []*review.Review
	for rows.Next() {
		var rev review.Review
		var images []byte
		err := rows.Scan(
			&rev.ID, &rev.ProjectID, &rev.OrderID, &rev.Rating, &rev.Comment, &images, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt, &rev.ClientName, &rev.ProjectName,
		)
		if err == nil {
			json.Unmarshal(images, &rev.Images)
			reviews = append(reviews, &rev)
		}
	}
	return reviews, total, nil
}

func (r *ReviewRepo) UpdateStatus(ctx context.Context, id int64, status review.ReviewStatus) error {
	_, err := r.pool.Exec(ctx, "UPDATE reviews SET status = $1, updated_at = NOW() WHERE id = $2", status, id)
	return err
}

func (r *ReviewRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM reviews WHERE id = $1", id)
	return err
}

func (r *ReviewRepo) GetByProjectID(ctx context.Context, projectID int64) ([]*review.Review, error) {
	query := `
		SELECT r.id, r.project_id, r.order_id, r.rating, r.comment, r.images, r.status, r.created_at, r.updated_at, o.client_name
		FROM reviews r
		JOIN orders o ON r.order_id = o.id
		WHERE r.project_id = $1 AND r.status = 'approved'
		ORDER BY r.created_at DESC`
	
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil { return nil, err }
	defer rows.Close()

	var reviews []*review.Review
	for rows.Next() {
		var rev review.Review
		var images []byte
		if err := rows.Scan(&rev.ID, &rev.ProjectID, &rev.OrderID, &rev.Rating, &rev.Comment, &images, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt, &rev.ClientName); err == nil {
			json.Unmarshal(images, &rev.Images)
			reviews = append(reviews, &rev)
		}
	}
	return reviews, nil
}
