package dto

import (
	"github.com/rostmebel/backend/internal/domain/review"
)

type CreateReviewRequest struct {
	ProjectID   *int64               `json:"project_id"`
	ClientPhone string               `json:"client_phone" validate:"required"`
	Rating      int                  `json:"rating" validate:"required,min=1,max=5"`
	Comment     string               `json:"comment" validate:"required"`
	Images      []review.ReviewImage `json:"images"`
}

type ReviewResponse struct {
	ID          int64                `json:"id"`
	ProjectID   *int64               `json:"project_id"`
	Rating      int                  `json:"rating"`
	Comment     string               `json:"comment"`
	Images      []review.ReviewImage `json:"images"`
	Status      review.ReviewStatus  `json:"status"`
	ClientName  string               `json:"client_name"`
	ProjectName string               `json:"project_name"`
	CreatedAt   string               `json:"created_at"`
}

func FromReview(r *review.Review) ReviewResponse {
	return ReviewResponse{
		ID:          r.ID,
		ProjectID:   r.ProjectID,
		Rating:      r.Rating,
		Comment:     r.Comment,
		Images:      r.Images,
		Status:      r.Status,
		ClientName:  r.ClientName,
		ProjectName: r.ProjectName,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type ModerateReviewRequest struct {
	Approved bool `json:"approved"`
}
