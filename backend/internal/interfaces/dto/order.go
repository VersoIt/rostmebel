package dto

import (
	"time"

	"github.com/rostmebel/backend/internal/domain/order"
)

type CreateOrderRequest struct {
	ProjectID     *int64 `json:"project_id"`
	ClientName    string `json:"client_name" validate:"required"`
	ClientPhone   string `json:"client_phone" validate:"required"`
	ClientEmail   string `json:"client_email" validate:"omitempty,email"`
	Comment       string `json:"comment" validate:"omitempty,max=2000"`
	ProjectType   string `json:"project_type" validate:"omitempty,max=80"`
	BudgetRange   string `json:"budget_range" validate:"omitempty,max=80"`
	City          string `json:"city" validate:"omitempty,max=120"`
	ContactMethod string `json:"contact_method" validate:"omitempty,oneof=phone whatsapp telegram max email"`
	Fingerprint   string `json:"fingerprint"`
	Website       string `json:"website"` // Honeypot
}

type OrderResponse struct {
	ID            int64             `json:"id"`
	ProjectID     *int64            `json:"project_id"`
	ProjectName   string            `json:"project_name,omitempty"`
	ClientName    string            `json:"client_name"`
	ClientPhone   string            `json:"client_phone"`
	ClientEmail   string            `json:"client_email"`
	Comment       string            `json:"comment"`
	ProjectType   string            `json:"project_type"`
	BudgetRange   string            `json:"budget_range"`
	City          string            `json:"city"`
	ContactMethod string            `json:"contact_method"`
	Status        order.OrderStatus `json:"status"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

func FromOrder(o *order.Order) OrderResponse {
	return OrderResponse{
		ID:            o.ID,
		ProjectID:     o.ProjectID,
		ProjectName:   o.ProjectName,
		ClientName:    o.ClientName,
		ClientPhone:   o.ClientPhone,
		ClientEmail:   o.ClientEmail,
		Comment:       o.Comment,
		ProjectType:   o.ProjectType,
		BudgetRange:   o.BudgetRange,
		City:          o.City,
		ContactMethod: o.ContactMethod,
		Status:        o.Status,
		CreatedAt:     o.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     o.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

type UpdateOrderStatusRequest struct {
	Status order.OrderStatus `json:"status" validate:"required"`
}
