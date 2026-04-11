package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rostmebel/backend/internal/domain/product"
)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*product.Project, error) {
	query := `SELECT id, project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status, views_count, orders_count, created_at, updated_at FROM projects WHERE id = $1 AND deleted_at IS NULL`
	var p product.Project
	var images []byte
	var specs []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.ProjectCategoryID, &p.Name, &p.Slug, &p.Description, &p.Budget, &p.BudgetOld, &images, &specs, &p.AITags, &p.Status, &p.ViewsCount, &p.OrdersCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("productRepo.GetByID: %w", err)
	}

	if err := json.Unmarshal(images, &p.Images); err != nil {
		return nil, fmt.Errorf("failed to unmarshal images: %w", err)
	}
	if err := json.Unmarshal(specs, &p.Details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal specs: %w", err)
	}

	return &p, nil
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*product.Project, error) {
	query := `SELECT id, project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status, views_count, orders_count, created_at, updated_at FROM projects WHERE slug = $1 AND deleted_at IS NULL`
	var p product.Project
	var images []byte
	var specs []byte

	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.ProjectCategoryID, &p.Name, &p.Slug, &p.Description, &p.Budget, &p.BudgetOld, &images, &specs, &p.AITags, &p.Status, &p.ViewsCount, &p.OrdersCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("productRepo.GetBySlug: %w", err)
	}

	if err := json.Unmarshal(images, &p.Images); err != nil {
		return nil, fmt.Errorf("failed to unmarshal images: %w", err)
	}
	if err := json.Unmarshal(specs, &p.Details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal specs: %w", err)
	}

	return &p, nil
}

func (r *ProductRepo) List(ctx context.Context, f product.ListFilter) ([]*product.Project, int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	conditions = append(conditions, "deleted_at IS NULL")

	if f.ProjectCategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("project_category_id = $%d", argCount))
		args = append(args, *f.ProjectCategoryID)
		argCount++
	}

	if f.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *f.Status)
		argCount++
	}

	if f.MinBudget != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argCount))
		args = append(args, *f.MinBudget)
		argCount++
	}

	if f.MaxBudget != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argCount))
		args = append(args, *f.MaxBudget)
		argCount++
	}

	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argCount))
		args = append(args, "%"+f.Search+"%")
		argCount++
	}

	if f.Cursor > 0 {
		conditions = append(conditions, fmt.Sprintf("id < $%d", argCount))
		args = append(args, f.Cursor)
		argCount++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects %s", where)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	order := "DESC"
	if strings.ToUpper(f.SortOrder) == "ASC" {
		order = "ASC"
	}
	sortBy := "id"
	if f.SortBy != "" {
		sortBy = f.SortBy
	}

	query := fmt.Sprintf(`
		SELECT id, project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status, views_count, orders_count, created_at, updated_at 
		FROM projects %s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`, 
		where, sortBy, order, argCount, argCount+1)
	
	args = append(args, f.Limit, f.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*product.Project
	for rows.Next() {
		var p product.Project
		var images []byte
		var specs []byte
		err := rows.Scan(
			&p.ID, &p.ProjectCategoryID, &p.Name, &p.Slug, &p.Description, &p.Budget, &p.BudgetOld, &images, &specs, &p.AITags, &p.Status, &p.ViewsCount, &p.OrdersCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan project: %w", err)
		}
		if err := json.Unmarshal(images, &p.Images); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal images: %w", err)
		}
		if err := json.Unmarshal(specs, &p.Details); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		projects = append(projects, &p)
	}

	return projects, total, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *product.Project) error {
	images, _ := json.Marshal(p.Images)
	specs, _ := json.Marshal(p.Details)

	query := `INSERT INTO projects (project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`
	
	return r.pool.QueryRow(ctx, query, 
		p.ProjectCategoryID, p.Name, p.Slug, p.Description, p.Budget, p.BudgetOld, images, specs, p.AITags, p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepo) Update(ctx context.Context, p *product.Project) error {
	images, _ := json.Marshal(p.Images)
	specs, _ := json.Marshal(p.Details)

	query := `UPDATE projects SET project_category_id = $1, name = $2, slug = $3, description = $4, price = $5, price_old = $6, images = $7, specs = $8, ai_tags = $9, status = $10, updated_at = NOW() 
			  WHERE id = $11 AND deleted_at IS NULL`
	
	_, err := r.pool.Exec(ctx, query, 
		p.ProjectCategoryID, p.Name, p.Slug, p.Description, p.Budget, p.BudgetOld, images, specs, p.AITags, p.Status, p.ID,
	)
	return err
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) ListCategories(ctx context.Context) ([]*product.Category, error) {
	query := `SELECT id, name, slug, icon, sort_order, created_at, updated_at FROM project_categories ORDER BY sort_order ASC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []*product.Category
	for rows.Next() {
		var c product.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *ProductRepo) GetCategoryBySlug(ctx context.Context, slug string) (*product.Category, error) {
	query := `SELECT id, name, slug, icon, sort_order, created_at, updated_at FROM project_categories WHERE slug = $1`
	var c product.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(&c.ID, &c.Name, &c.Slug, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ProductRepo) CreateCategory(ctx context.Context, c *product.Category) error {
	query := `INSERT INTO project_categories (name, slug, icon, sort_order) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, c.Name, c.Slug, c.Icon, c.SortOrder).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ProductRepo) UpdateCategory(ctx context.Context, c *product.Category) error {
	query := `UPDATE project_categories SET name = $1, slug = $2, icon = $3, sort_order = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.pool.Exec(ctx, query, c.Name, c.Slug, c.Icon, c.SortOrder, c.ID)
	return err
}

func (r *ProductRepo) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM project_categories WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) IncrementViews(ctx context.Context, id int64) error {
	query := `UPDATE projects SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) Search(ctx context.Context, query string, limit int) ([]*product.Project, error) {
	sql := `
		SELECT id, project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status, views_count, orders_count, created_at, updated_at 
		FROM projects 
		WHERE deleted_at IS NULL AND status = 'published' AND search_vector @@ plainto_tsquery('russian', $1)
		ORDER BY ts_rank(search_vector, plainto_tsquery('russian', $1)) DESC
		LIMIT $2`
	
	rows, err := r.pool.Query(ctx, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}
	defer rows.Close()

	var projects []*product.Project
	for rows.Next() {
		var p product.Project
		var images []byte
		var specs []byte
		err := rows.Scan(
			&p.ID, &p.ProjectCategoryID, &p.Name, &p.Slug, &p.Description, &p.Budget, &p.BudgetOld, &images, &specs, &p.AITags, &p.Status, &p.ViewsCount, &p.OrdersCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		json.Unmarshal(images, &p.Images)
		json.Unmarshal(specs, &p.Details)
		projects = append(projects, &p)
	}
	return projects, nil
}
