package postgres

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rostmebel/backend/internal/domain/order"
)

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

func (r *OrderRepo) Create(ctx context.Context, o *order.Order) error {
	query := `
		INSERT INTO orders (
			project_id, client_name, client_phone, client_email, comment,
			project_type, budget_range, city, contact_method,
			status, ip_address, user_agent, fingerprint
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		o.ProjectID, o.ClientName, o.ClientPhone, o.ClientEmail, o.Comment,
		o.ProjectType, o.BudgetRange, o.City, o.ContactMethod,
		o.Status, o.IPAddress.String(), o.UserAgent, o.Fingerprint,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	query := `
		SELECT o.id, o.project_id, o.client_name, o.client_phone,
			COALESCE(o.client_email, ''), COALESCE(o.comment, ''),
			COALESCE(o.project_type, ''), COALESCE(o.budget_range, ''),
			COALESCE(o.city, ''), COALESCE(o.contact_method, ''),
			o.status, o.ip_address::text, COALESCE(o.user_agent, ''), COALESCE(o.fingerprint, ''),
			o.created_at, o.updated_at, p.name
		FROM orders o
		LEFT JOIN projects p ON o.project_id = p.id
		WHERE o.id = $1`

	var o order.Order
	var ip string
	var projectName *string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone,
		&o.ClientEmail, &o.Comment, &o.ProjectType, &o.BudgetRange, &o.City, &o.ContactMethod,
		&o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt, &projectName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if projectName != nil {
		o.ProjectName = *projectName
	}
	o.IPAddress = net.ParseIP(ip)
	return &o, nil
}

// List returns filtered orders, filtered count and absolute total count
func (r *OrderRepo) List(ctx context.Context, f order.ListFilter) ([]*order.Order, int64, int64, error) {
	where := ""
	args := []interface{}{}
	if f.Status != "" {
		where = "WHERE o.status = $1"
		args = append(args, f.Status)
	}

	// Filtered total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders o %s", where)
	var filteredTotal int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&filteredTotal); err != nil {
		return nil, 0, 0, err
	}

	// Absolute total
	var absoluteTotal int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&absoluteTotal); err != nil {
		return nil, 0, 0, err
	}

	limit := 20
	if f.Limit > 0 {
		limit = f.Limit
	}
	offset := 0
	if f.Offset > 0 {
		offset = f.Offset
	}

	query := fmt.Sprintf(`
		SELECT o.id, o.project_id, o.client_name, o.client_phone,
			COALESCE(o.client_email, ''), COALESCE(o.comment, ''),
			COALESCE(o.project_type, ''), COALESCE(o.budget_range, ''),
			COALESCE(o.city, ''), COALESCE(o.contact_method, ''),
			o.status, o.ip_address::text, COALESCE(o.user_agent, ''), COALESCE(o.fingerprint, ''),
			o.created_at, o.updated_at, p.name
		FROM orders o
		LEFT JOIN projects p ON o.project_id = p.id
		%s 
		ORDER BY o.created_at DESC 
		LIMIT $%d OFFSET $%d`, where, len(args)+1, len(args)+2)

	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var ip string
		var projectName *string
		err := rows.Scan(
			&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone,
			&o.ClientEmail, &o.Comment, &o.ProjectType, &o.BudgetRange, &o.City, &o.ContactMethod,
			&o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt, &projectName,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		if projectName != nil {
			o.ProjectName = *projectName
		}
		o.IPAddress = net.ParseIP(ip)
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}
	return orders, filteredTotal, absoluteTotal, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id int64, status order.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *OrderRepo) GetOrderCountByIP(ctx context.Context, ip net.IP, since time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM orders WHERE ip_address = $1 AND created_at >= $2`
	var count int
	err := r.pool.QueryRow(ctx, query, ip.String(), since).Scan(&count)
	return count, err
}

func (r *OrderRepo) IsIPBlocked(ctx context.Context, ip net.IP) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM ip_blocks WHERE ip_address = $1 AND expires_at > NOW())`
	var blocked bool
	err := r.pool.QueryRow(ctx, query, ip.String()).Scan(&blocked)
	return blocked, err
}

func (r *OrderRepo) MarkAsSpam(ctx context.Context, id int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var ip string
	err = tx.QueryRow(ctx, "UPDATE orders SET status = 'spam', updated_at = NOW() WHERE id = $1 RETURNING ip_address::text", id).Scan(&ip)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO ip_blocks (ip_address, reason, expires_at) VALUES ($1, 'spam', $2) ON CONFLICT (ip_address) DO UPDATE SET expires_at = $2", ip, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrderRepo) UnblockIP(ctx context.Context, ip net.IP) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM ip_blocks WHERE ip_address = $1", ip.String())
	return err
}

func (r *OrderRepo) Export(ctx context.Context) ([]*order.Order, error) {
	query := `
		SELECT o.id, o.project_id, o.client_name, o.client_phone,
			COALESCE(o.client_email, ''), COALESCE(o.comment, ''),
			COALESCE(o.project_type, ''), COALESCE(o.budget_range, ''),
			COALESCE(o.city, ''), COALESCE(o.contact_method, ''),
			o.status, o.ip_address::text, COALESCE(o.user_agent, ''), COALESCE(o.fingerprint, ''),
			o.created_at, o.updated_at, p.name
		FROM orders o
		LEFT JOIN projects p ON o.project_id = p.id
		ORDER BY o.created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var ip string
		var projectName *string
		err := rows.Scan(
			&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone,
			&o.ClientEmail, &o.Comment, &o.ProjectType, &o.BudgetRange, &o.City, &o.ContactMethod,
			&o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt, &projectName,
		)
		if err != nil {
			return nil, err
		}
		if projectName != nil {
			o.ProjectName = *projectName
		}
		o.IPAddress = net.ParseIP(ip)
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
