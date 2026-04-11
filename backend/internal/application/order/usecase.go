package order

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/order"
)

type UseCase struct {
	repo  order.Repository
	redis *redis.Client
}

func NewUseCase(repo order.Repository, redis *redis.Client) *UseCase {
	return &UseCase{repo: repo, redis: redis}
}

func (u *UseCase) CreateOrder(ctx context.Context, o *order.Order) error {
	// Rate limiting: 3 orders per IP per 24 hours
	key := fmt.Sprintf("order_limit:%s", o.IPAddress.String())
	count, err := u.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if count >= 3 {
		return fmt.Errorf("rate limit exceeded: max 3 orders per 24 hours")
	}

	if err := u.repo.Create(ctx, o); err != nil {
		return err
	}

	// Increment rate limit
	if err := u.redis.Incr(ctx, key).Err(); err != nil {
		return err
	}
	if count == 0 {
		u.redis.Expire(ctx, key, 24*time.Hour)
	}

	return nil
}

func (u *UseCase) GetOrder(ctx context.Context, id int64) (*order.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) ListOrders(ctx context.Context, f order.ListFilter) ([]*order.Order, int64, error) {
	return u.repo.List(ctx, f)
}

func (u *UseCase) UpdateOrderStatus(ctx context.Context, id int64, status order.OrderStatus) error {
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *UseCase) MarkAsSpam(ctx context.Context, id int64) error {
	return u.repo.MarkAsSpam(ctx, id)
}

func (u *UseCase) ExportOrders(ctx context.Context) ([]*order.Order, error) {
	return u.repo.Export(ctx)
}
