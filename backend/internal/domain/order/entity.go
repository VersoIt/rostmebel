package order

import (
	"context"
	"net"
	"time"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "new"
	StatusProcessing OrderStatus = "processing"
	StatusDone       OrderStatus = "done"
	StatusRejected   OrderStatus = "rejected"
	StatusSpam       OrderStatus = "spam"
)

type Order struct {
	ID          int64       `json:"id"`
	ProjectID   *int64      `json:"project_id"`
	ProjectName string      `json:"project_name,omitempty"` // Added for admin view
	ClientName  string      `json:"client_name"`
	ClientPhone string      `json:"client_phone"`
	ClientEmail string      `json:"client_email"`
	Comment     string      `json:"comment"`
	Status      OrderStatus `json:"status"`
	IPAddress   net.IP      `json:"ip_address"`
	UserAgent   string      `json:"user_agent"`
	Fingerprint string      `json:"fingerprint"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	List(ctx context.Context, filter ListFilter) ([]*Order, int64, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus) error
	GetOrderCountByIP(ctx context.Context, ip net.IP, since time.Time) (int, error)
	MarkAsSpam(ctx context.Context, id int64) error
	Export(ctx context.Context) ([]*Order, error)
}

type ListFilter struct {
	Status OrderStatus
	Limit  int
	Offset int
}
