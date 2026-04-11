package admin

import (
	"context"
	"time"
)

type Admin struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	RefreshToken *string    `json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*Admin, error)
	GetByID(ctx context.Context, id int64) (*Admin, error)
	Create(ctx context.Context, a *Admin) error
	Update(ctx context.Context, a *Admin) error
	UpdateRefreshToken(ctx context.Context, id int64, token *string) error
	GetStats(ctx context.Context) (*Stats, error)
}

type Stats struct {
	ProductsCount  int64            `json:"products_count"`
	NewOrdersToday int64            `json:"new_orders_today"`
	TotalOrders    int64            `json:"total_orders"`
	SuccessRate    float64          `json:"success_rate"`
	TopProducts    []TopProduct     `json:"top_products"`
	OrdersByDay    []OrdersByDay    `json:"orders_by_day"`
}

type TopProduct struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type OrdersByDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
