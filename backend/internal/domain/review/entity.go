package review

import (
	"context"
	"time"
)

type ReviewStatus string

const (
	StatusPending  ReviewStatus = "pending"
	StatusApproved ReviewStatus = "approved"
	StatusRejected ReviewStatus = "rejected"
)

type ReviewImage struct {
	URL string `json:"url"`
}

type Review struct {
	ID          int64         `json:"id"`
	ProjectID   *int64        `json:"project_id"`
	OrderID     int64         `json:"order_id"`
	Rating      int           `json:"rating"`
	Comment     string        `json:"comment"`
	Images      []ReviewImage `json:"images"`
	Status      ReviewStatus  `json:"status"`
	
	// Join fields for frontend
	ClientName  string `json:"client_name,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, r *Review) error
	GetByID(ctx context.Context, id int64) (*Review, error)
	List(ctx context.Context, filter ListFilter) ([]*Review, int64, error)
	UpdateStatus(ctx context.Context, id int64, status ReviewStatus) error
	Delete(ctx context.Context, id int64) error
	GetByProjectID(ctx context.Context, projectID int64) ([]*Review, error)
}

type ListFilter struct {
	Status    ReviewStatus
	ProjectID *int64
	Limit     int
	Offset    int
}
