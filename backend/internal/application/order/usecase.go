package order

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/order"
	domProduct "github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/infrastructure/telegram"
)

type UseCase struct {
	repo         order.Repository
	prodRepo     domProduct.Repository
	redis        *redis.Client
	tg           *telegram.Client
	limitEnabled bool
}

func NewUseCase(repo order.Repository, prodRepo domProduct.Repository, redis *redis.Client, tg *telegram.Client, limitEnabled bool) *UseCase {
	return &UseCase{repo: repo, prodRepo: prodRepo, redis: redis, tg: tg, limitEnabled: limitEnabled}
}

func (u *UseCase) CreateOrder(ctx context.Context, o *order.Order) error {
	if u.limitEnabled {
		// Rate limiting: 1 order per IP per 24 hours
		key := fmt.Sprintf("order_limit:%s", o.IPAddress.String())
		count, err := u.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if count >= 1 {
			return fmt.Errorf("Вы уже отправляли заявку сегодня. Пожалуйста, подождите 24 часа.")
		}
		
		// Increment rate limit later
		defer func() {
			u.redis.Incr(ctx, key)
			u.redis.Expire(ctx, key, 24*time.Hour)
		}()
	}

	if err := u.repo.Create(ctx, o); err != nil {
		return err
	}

	// Telegram Notification
	productName := "Общая консультация"
	if o.ProductID != nil {
		if p, err := u.prodRepo.GetByID(ctx, *o.ProductID); err == nil && p != nil {
			productName = p.Name
		}
	}
	go u.tg.SendOrderNotification(o.ClientName, o.ClientPhone, productName, o.Comment)

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
