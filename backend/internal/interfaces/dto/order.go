package dto

import (
	"github.com/rostmebel/backend/internal/domain/order"
)

type CreateOrderRequest struct {
	ProductID   *int64  `json:"product_id"`
	ClientName  string  `json:"client_name" validate:"required"`
	ClientPhone string  `json:"client_phone" validate:"required"`
	ClientEmail string  `json:"client_email"`
	Comment     string  `json:"comment"`
	Fingerprint string  `json:"fingerprint"`
	Website     string  `json:"website"` // Honeypot
}

type OrderResponse struct {
	ID          int64             `json:"id"`
	ProductID   *int64            `json:"product_id"`
	ClientName  string            `json:"client_name"`
	ClientPhone string            `json:"client_phone"`
	ClientEmail string            `json:"client_email"`
	Comment     string            `json:"comment"`
	Status      order.OrderStatus `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

func FromOrder(o *order.Order) OrderResponse {
	return OrderResponse{
		ID:          o.ID,
		ProductID:   o.ProductID,
		ClientName:  o.ClientName,
		ClientPhone: o.ClientPhone,
		ClientEmail: o.ClientEmail,
		Comment:     o.Comment,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type UpdateOrderStatusRequest struct {
	Status order.OrderStatus `json:"status" validate:"required"`
}
