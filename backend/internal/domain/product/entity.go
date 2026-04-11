package product

import (
	"context"
	"time"
)

type ProductStatus string

const (
	StatusPublished ProductStatus = "published"
	StatusDraft     ProductStatus = "draft"
	StatusArchived  ProductStatus = "archived"
)

type Image struct {
	URL    string `json:"url"`
	IsMain bool   `json:"is_main"`
}

type Product struct {
	ID           int64             `json:"id"`
	CategoryID   *int64            `json:"category_id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	Description  string            `json:"description"`
	Price        float64           `json:"price"`
	PriceOld     *float64          `json:"price_old"`
	Images       []Image           `json:"images"`
	Specs        map[string]string `json:"specs"`
	AITags       string            `json:"ai_tags"`
	Status       ProductStatus     `json:"status"`
	ViewsCount   int               `json:"views_count"`
	OrdersCount  int               `json:"orders_count"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    *time.Time        `json:"deleted_at,omitempty"`
}

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	List(ctx context.Context, filter ListFilter) ([]*Product, int64, error)
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
	
	ListCategories(ctx context.Context) ([]*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id int64) error

	IncrementViews(ctx context.Context, id int64) error
	Search(ctx context.Context, query string, limit int) ([]*Product, error)
}

type ListFilter struct {
	CategoryID *int64
	Status     *ProductStatus
	MinPrice   *float64
	MaxPrice   *float64
	Search     string
	SortBy     string
	SortOrder  string
	Limit      int
	Offset     int
	Cursor     int64
}
