package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rostmebel/backend/internal/domain/admin"
)

type AdminRepo struct {
	pool *pgxpool.Pool
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{pool: pool}
}

func (r *AdminRepo) GetByUsername(ctx context.Context, username string) (*admin.Admin, error) {
	query := `SELECT id, username, password_hash, refresh_token, last_login_at, created_at FROM admins WHERE username = $1`
	var a admin.Admin
	err := r.pool.QueryRow(ctx, query, username).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.RefreshToken, &a.LastLoginAt, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) GetByID(ctx context.Context, id int64) (*admin.Admin, error) {
	query := `SELECT id, username, password_hash, refresh_token, last_login_at, created_at FROM admins WHERE id = $1`
	var a admin.Admin
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.RefreshToken, &a.LastLoginAt, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) Create(ctx context.Context, a *admin.Admin) error {
	query := `INSERT INTO admins (username, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, a.Username, a.PasswordHash).Scan(&a.ID, &a.CreatedAt)
}

func (r *AdminRepo) Update(ctx context.Context, a *admin.Admin) error {
	query := `UPDATE admins SET username = $1, password_hash = $2, last_login_at = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, a.Username, a.PasswordHash, a.LastLoginAt, a.ID)
	return err
}

func (r *AdminRepo) UpdateRefreshToken(ctx context.Context, id int64, token *string) error {
	query := `UPDATE admins SET refresh_token = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, token, id)
	return err
}

func (r *AdminRepo) GetStats(ctx context.Context) (*admin.Stats, error) {
	var stats admin.Stats

	// Products count
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL").Scan(&stats.ProductsCount)
	if err != nil {
		return nil, err
	}

	// New orders today
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE").Scan(&stats.NewOrdersToday)
	if err != nil {
		return nil, err
	}

	// Total orders
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	if err != nil {
		return nil, err
	}

	// Success rate (done orders / total orders)
	if stats.TotalOrders > 0 {
		var doneOrders int64
		r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE status = 'done'").Scan(&doneOrders)
		stats.SuccessRate = float64(doneOrders) / float64(stats.TotalOrders) * 100
	} else {
		stats.SuccessRate = 0
	}

	// Top products
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.name, COUNT(o.id) as count 
		FROM products p 
		LEFT JOIN orders o ON p.id = o.product_id 
		WHERE p.deleted_at IS NULL 
		GROUP BY p.id 
		ORDER BY count DESC 
		LIMIT 3`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp admin.TopProduct
			if err := rows.Scan(&tp.ID, &tp.Name, &tp.Count); err == nil {
				stats.TopProducts = append(stats.TopProducts, tp)
			}
		}
	}

	// Orders by day (last 30 days)
	rows, err = r.pool.Query(ctx, `
		SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as date, COUNT(*) 
		FROM orders 
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days' 
		GROUP BY date 
		ORDER BY date ASC`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var obd admin.OrdersByDay
			if err := rows.Scan(&obd.Date, &obd.Count); err == nil {
				stats.OrdersByDay = append(stats.OrdersByDay, obd)
			}
		}
	}

	return &stats, nil
}
