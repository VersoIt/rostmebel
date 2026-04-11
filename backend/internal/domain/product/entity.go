package product

import (
	"context"
	"time"
)

type ProjectStatus string

const (
	StatusPublished ProjectStatus = "published"
	StatusDraft     ProjectStatus = "draft"
	StatusArchived  ProjectStatus = "archived"
)

type Image struct {
	URL    string `json:"url"`
	IsMain bool   `json:"is_main"`
}

type Project struct {
	ID                int64             `json:"id"`
	ProjectCategoryID *int64            `json:"project_category_id"`
	Name              string            `json:"name"`
	Slug              string            `json:"slug"`
	Description       string            `json:"description"`
	Budget            float64           `json:"price"` // mapped to 'price' in DB for compatibility if needed, but here let's call it Budget
	BudgetOld         *float64          `json:"price_old"`
	Images            []Image           `json:"images"`
	Details           map[string]string `json:"specs"` // renamed from specs
	AITags            string            `json:"ai_tags"`
	Status            ProjectStatus     `json:"status"`
	ViewsCount        int               `json:"views_count"`
	OrdersCount       int               `json:"orders_count"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
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
	GetByID(ctx context.Context, id int64) (*Project, error)
	GetBySlug(ctx context.Context, slug string) (*Project, error)
	List(ctx context.Context, filter ListFilter) ([]*Project, int64, error)
	Create(ctx context.Context, p *Project) error
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, id int64) error
	
	ListCategories(ctx context.Context) ([]*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id int64) error

	IncrementViews(ctx context.Context, id int64) error
	Search(ctx context.Context, query string, limit int) ([]*Project, error)
}

type ListFilter struct {
	ProjectCategoryID *int64
	Status            *ProjectStatus
	MinBudget         *float64
	MaxBudget         *float64
	Search            string
	SortBy            string
	SortOrder         string
	Limit             int
	Offset            int
	Cursor            int64
}
