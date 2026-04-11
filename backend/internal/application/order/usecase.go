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
		// Rate limiting: 5 orders per IP per 24 hours (increased from 1)
		key := fmt.Sprintf("order_limit:%s", o.IPAddress.String())
		count, err := u.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if count >= 5 {
			return fmt.Errorf("Превышен лимит заявок на сегодня. Пожалуйста, попробуйте завтра.")
		}
		
		// Increment rate limit later
		defer func() {
			u.redis.Incr(ctx, key)
			u.redis.Expire(ctx, key, 24*time.Hour)
		}()
	}

	if err := u.repo.Create(ctx, o); err != nil {
		return fmt.Errorf("failed to save order to DB: %w", err)
	}

	// Telegram Notification
	projectName := "Общая консультация"
	if o.ProjectID != nil {
		if p, err := u.prodRepo.GetByID(ctx, *o.ProjectID); err == nil && p != nil {
			projectName = p.Name
		}
	}
	
	// Use background context for notification to avoid cancelling it if user request finishes
	go func() {
		err := u.tg.SendOrderNotification(o.ClientName, o.ClientPhone, projectName, o.Comment)
		if err != nil {
			fmt.Printf("TELEGRAM ERROR: %v\n", err)
		}
	}()

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
