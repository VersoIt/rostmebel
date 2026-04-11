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
	query := `INSERT INTO orders (project_id, client_name, client_phone, client_email, comment, status, ip_address, user_agent, fingerprint) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	
	return r.pool.QueryRow(ctx, query, 
		o.ProjectID, o.ClientName, o.ClientPhone, o.ClientEmail, o.Comment, o.Status, o.IPAddress.String(), o.UserAgent, o.Fingerprint,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	query := `SELECT id, project_id, client_name, client_phone, client_email, comment, status, ip_address, user_agent, fingerprint, created_at, updated_at FROM orders WHERE id = $1`
	var o order.Order
	var ip string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone, &o.ClientEmail, &o.Comment, &o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	o.IPAddress = net.ParseIP(ip)
	return &o, nil
}

func (r *OrderRepo) List(ctx context.Context, f order.ListFilter) ([]*order.Order, int64, error) {
	where := ""
	args := []interface{}{}
	if f.Status != "" {
		where = "WHERE status = $1"
		args = append(args, f.Status)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", where)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := 20
	if f.Limit > 0 {
		limit = f.Limit
	}
	offset := 0
	if f.Offset > 0 {
		offset = f.Offset
	}

	query := fmt.Sprintf("SELECT id, project_id, client_name, client_phone, client_email, comment, status, ip_address, user_agent, fingerprint, created_at, updated_at FROM orders %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", where, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var ip string
		if err := rows.Scan(&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone, &o.ClientEmail, &o.Comment, &o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		o.IPAddress = net.ParseIP(ip)
		orders = append(orders, &o)
	}
	return orders, total, nil
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

func (r *OrderRepo) MarkAsSpam(ctx context.Context, id int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var ip string
	err = tx.QueryRow(ctx, "UPDATE orders SET status = 'spam' WHERE id = $1 RETURNING ip_address", id).Scan(&ip)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO ip_blocks (ip_address, reason, expires_at) VALUES ($1, 'spam', $2) ON CONFLICT (ip_address) DO UPDATE SET expires_at = $2", ip, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrderRepo) Export(ctx context.Context) ([]*order.Order, error) {
	query := `SELECT id, project_id, client_name, client_phone, client_email, comment, status, ip_address, user_agent, fingerprint, created_at, updated_at FROM orders ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var ip string
		if err := rows.Scan(&o.ID, &o.ProjectID, &o.ClientName, &o.ClientPhone, &o.ClientEmail, &o.Comment, &o.Status, &ip, &o.UserAgent, &o.Fingerprint, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.IPAddress = net.ParseIP(ip)
		orders = append(orders, &o)
	}
	return orders, nil
}
