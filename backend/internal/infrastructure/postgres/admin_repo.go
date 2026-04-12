package postgres

import (
	"context"
	"fmt"

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
	stats := &admin.Stats{
		TopProjects:  []admin.TopProject{},
		OrdersByDay:  []admin.OrdersByDay{},
		RecentOrders: []admin.RecentOrder{},
	}

	// 1. Projects count
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL").Scan(&stats.ProjectsCount)
	if err != nil {
		return nil, fmt.Errorf("stats projects count: %w", err)
	}

	// 2. New orders today
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE AND status != 'spam'").Scan(&stats.NewOrdersToday)
	if err != nil {
		return nil, fmt.Errorf("stats new orders: %w", err)
	}

	// 3. Total orders
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE status != 'spam'").Scan(&stats.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("stats total orders: %w", err)
	}

	// 4. Success rate
	var doneOrders int64
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE status = 'done'").Scan(&doneOrders)
	if err != nil {
		return nil, fmt.Errorf("stats success rate: %w", err)
	}
	if stats.TotalOrders > 0 {
		stats.SuccessRate = float64(doneOrders) / float64(stats.TotalOrders) * 100
	}

	// 5. Top projects
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.name, COUNT(o.id) as count
		FROM projects p
		JOIN orders o ON o.project_id = p.id
		WHERE o.status != 'spam' AND p.deleted_at IS NULL
		GROUP BY p.id, p.name
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, fmt.Errorf("stats top projects: %w", err)
	}
	for rows.Next() {
		var tp admin.TopProject
		if err := rows.Scan(&tp.ID, &tp.Name, &tp.Count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("stats top projects scan: %w", err)
		}
		stats.TopProjects = append(stats.TopProjects, tp)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("stats top projects rows: %w", err)
	}
	rows.Close()

	// 6. Recent Orders (LATEST 5)
	rows, err = r.pool.Query(ctx, `
		SELECT o.id, o.client_name, COALESCE(p.name, 'Общая консультация'), o.created_at
		FROM orders o
		LEFT JOIN projects p ON o.project_id = p.id
		WHERE o.status = 'new'
		ORDER BY o.created_at DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, fmt.Errorf("stats recent orders: %w", err)
	}
	for rows.Next() {
		var ro admin.RecentOrder
		if err := rows.Scan(&ro.ID, &ro.ClientName, &ro.ProjectName, &ro.CreatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("stats recent orders scan: %w", err)
		}
		stats.RecentOrders = append(stats.RecentOrders, ro)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("stats recent orders rows: %w", err)
	}
	rows.Close()

	// 7. Pending Reviews Count
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM reviews WHERE status = 'pending'").Scan(&stats.PendingReviewsCount)
	if err != nil {
		return nil, fmt.Errorf("stats pending reviews: %w", err)
	}

	// 8. Orders by day
	rows, err = r.pool.Query(ctx, `
		SELECT TO_CHAR(d, 'YYYY-MM-DD') as date, COUNT(o.id) as count
		FROM (SELECT generate_series(CURRENT_DATE - INTERVAL '29 days', CURRENT_DATE, '1 day')::date as d) d
		LEFT JOIN orders o ON DATE(o.created_at) = d AND o.status != 'spam'
		GROUP BY d
		ORDER BY d ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("stats orders by day: %w", err)
	}
	for rows.Next() {
		var obd admin.OrdersByDay
		if err := rows.Scan(&obd.Date, &obd.Count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("stats orders by day scan: %w", err)
		}
		stats.OrdersByDay = append(stats.OrdersByDay, obd)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("stats orders by day rows: %w", err)
	}
	rows.Close()

	return stats, nil
}
