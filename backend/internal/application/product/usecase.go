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

func (u *UseCase) GetProject(ctx context.Context, id int64) (*product.Project, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) GetProjectBySlug(ctx context.Context, slug string) (*product.Project, error) {
	return u.repo.GetBySlug(ctx, slug)
}

func (u *UseCase) ListProjects(ctx context.Context, f product.ListFilter) ([]*product.Project, int64, error) {
	return u.repo.List(ctx, f)
}

func (u *UseCase) CreateProject(ctx context.Context, p *product.Project) error {
	return u.repo.Create(ctx, p)
}

func (u *UseCase) UpdateProject(ctx context.Context, p *product.Project) error {
	return u.repo.Update(ctx, p)
}

func (u *UseCase) DeleteProject(ctx context.Context, id int64) error {
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

func (u *UseCase) Search(ctx context.Context, query string, limit int) ([]*product.Project, error) {
	return u.repo.Search(ctx, query, limit)
}
