package product

import (
	"context"

	"github.com/rostmebel/backend/internal/domain/product"
)

type UseCase struct {
	repo product.Repository
}

func NewUseCase(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetProduct(ctx context.Context, id int64) (*product.Product, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	return u.repo.GetBySlug(ctx, slug)
}

func (u *UseCase) ListProducts(ctx context.Context, f product.ListFilter) ([]*product.Product, int64, error) {
	return u.repo.List(ctx, f)
}

func (u *UseCase) CreateProduct(ctx context.Context, p *product.Product) error {
	return u.repo.Create(ctx, p)
}

func (u *UseCase) UpdateProduct(ctx context.Context, p *product.Product) error {
	return u.repo.Update(ctx, p)
}

func (u *UseCase) DeleteProduct(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *UseCase) ListCategories(ctx context.Context) ([]*product.Category, error) {
	return u.repo.ListCategories(ctx)
}

func (u *UseCase) CreateCategory(ctx context.Context, c *product.Category) error {
	return u.repo.CreateCategory(ctx, c)
}

func (u *UseCase) UpdateCategory(ctx context.Context, c *product.Category) error {
	return u.repo.UpdateCategory(ctx, c)
}

func (u *UseCase) DeleteCategory(ctx context.Context, id int64) error {
	return u.repo.DeleteCategory(ctx, id)
}

func (u *UseCase) IncrementViews(ctx context.Context, id int64) error {
	return u.repo.IncrementViews(ctx, id)
}

func (u *UseCase) Search(ctx context.Context, query string, limit int) ([]*product.Product, error) {
	return u.repo.Search(ctx, query, limit)
}
