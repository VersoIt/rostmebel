package dto

import (
	"github.com/rostmebel/backend/internal/domain/product"
)

type ProductResponse struct {
	ID          int64                    `json:"id"`
	CategoryID  *int64                   `json:"category_id"`
	Name        string                   `json:"name"`
	Slug        string                   `json:"slug"`
	Description string                   `json:"description"`
	Price       float64                  `json:"price"`
	PriceOld    *float64                 `json:"price_old"`
	Images      []product.Image          `json:"images"`
	Specs       map[string]string        `json:"specs"`
	AITags      string                   `json:"ai_tags"`
	Status      product.ProductStatus    `json:"status"`
	ViewsCount  int                      `json:"views_count"`
	OrdersCount int                      `json:"orders_count"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

func FromProduct(p *product.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		PriceOld:    p.PriceOld,
		Images:      p.Images,
		Specs:       p.Specs,
		AITags:      p.AITags,
		Status:      p.Status,
		ViewsCount:  p.ViewsCount,
		OrdersCount: p.OrdersCount,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type CreateProductRequest struct {
	CategoryID  *int64                `json:"category_id"`
	Name        string                `json:"name" validate:"required"`
	Slug        string                `json:"slug" validate:"required"`
	Description string                `json:"description"`
	Price       float64               `json:"price" validate:"required"`
	PriceOld    *float64              `json:"price_old"`
	Images      []product.Image       `json:"images"`
	Specs       map[string]string     `json:"specs"`
	AITags      string                `json:"ai_tags"`
	Status      product.ProductStatus `json:"status"`
}

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

func FromCategory(c *product.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		Icon:      c.Icon,
		SortOrder: c.SortOrder,
	}
}

type AISearchRequest struct {
	Query string `json:"query" validate:"required"`
}
