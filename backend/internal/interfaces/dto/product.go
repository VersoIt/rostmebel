package dto

import (
	"github.com/rostmebel/backend/internal/domain/product"
)

type ProjectResponse struct {
	ID                int64                    `json:"id"`
	ProjectCategoryID *int64                   `json:"project_category_id"`
	Name              string                   `json:"name"`
	Slug              string                   `json:"slug"`
	Description       string                   `json:"description"`
	Budget            float64                  `json:"price"`
	BudgetOld         *float64                 `json:"price_old"`
	Images            []product.Image          `json:"images"`
	Details           map[string]string        `json:"specs"`
	AITags            string                   `json:"ai_tags"`
	Status            product.ProjectStatus    `json:"status"`
	ViewsCount        int                      `json:"views_count"`
	OrdersCount       int                      `json:"orders_count"`
	CreatedAt         string                   `json:"created_at"`
	UpdatedAt         string                   `json:"updated_at"`
}

func FromProject(p *product.Project) ProjectResponse {
	return ProjectResponse{
		ID:                p.ID,
		ProjectCategoryID: p.ProjectCategoryID,
		Name:              p.Name,
		Slug:              p.Slug,
		Description:       p.Description,
		Budget:            p.Budget,
		BudgetOld:         p.BudgetOld,
		Images:            p.Images,
		Details:           p.Details,
		AITags:            p.AITags,
		Status:            p.Status,
		ViewsCount:        p.ViewsCount,
		OrdersCount:       p.OrdersCount,
		CreatedAt:         p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type CreateProjectRequest struct {
	ProjectCategoryID *int64                `json:"project_category_id"`
	Name              string                `json:"name" validate:"required"`
	Slug              string                `json:"slug" validate:"required"`
	Description       string                `json:"description"`
	Budget            float64               `json:"price" validate:"required"`
	BudgetOld         *float64              `json:"price_old"`
	Images            []product.Image       `json:"images"`
	Details           map[string]string     `json:"specs"`
	AITags            string                `json:"ai_tags"`
	Status            product.ProjectStatus `json:"status"`
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
