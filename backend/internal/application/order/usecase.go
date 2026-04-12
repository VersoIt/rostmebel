package order

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/apperror"
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
	// 1. Check if IP is blocked
	blocked, err := u.repo.IsIPBlocked(ctx, o.IPAddress)
	if err != nil {
		return err
	}
	if blocked {
		return apperror.New(apperror.CodeOrderIPBlocked, "Client IP is temporarily blocked", map[string]any{
			"reason": "spam",
		})
	}

	if u.limitEnabled {
		// Rate limiting: 5 orders per IP per 24 hours
		key := fmt.Sprintf("order_limit:%s", o.IPAddress.String())
		count, err := u.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if count >= 5 {
			return apperror.New(apperror.CodeOrderRateLimited, "Order rate limit exceeded", map[string]any{
				"limit":        5,
				"window_hours": 24,
			})
		}

		defer func() {
			u.redis.Incr(ctx, key)
			u.redis.Expire(ctx, key, 24*time.Hour)
		}()
	}

	if err := u.repo.Create(ctx, o); err != nil {
		return apperror.Wrap(err, apperror.CodeInternal, "Failed to create order", nil)
	}

	// Telegram Notification
	projectName := "Общая консультация"
	if o.ProjectID != nil {
		if p, err := u.prodRepo.GetByID(ctx, *o.ProjectID); err == nil && p != nil {
			projectName = p.Name
		}
	}

	go func() {
		err := u.tg.SendOrderNotification(telegram.OrderNotification{
			Name:          o.ClientName,
			Phone:         o.ClientPhone,
			Product:       projectName,
			Comment:       o.Comment,
			ProjectType:   o.ProjectType,
			BudgetRange:   o.BudgetRange,
			City:          o.City,
			ContactMethod: o.ContactMethod,
		})
		if err != nil {
			fmt.Printf("TELEGRAM ERROR: %v\n", err)
		}
	}()

	return nil
}

func (u *UseCase) GetOrder(ctx context.Context, id int64) (*order.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) ListOrders(ctx context.Context, f order.ListFilter) ([]*order.Order, int64, int64, error) {
	return u.repo.List(ctx, f)
}

func (u *UseCase) UpdateOrderStatus(ctx context.Context, id int64, status order.OrderStatus) error {
	oldOrder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if oldOrder == nil {
		return apperror.New(apperror.CodeOrderNotFound, "Order not found", map[string]any{
			"id": id,
		})
	}

	// If restoring from spam, unblock IP.
	if oldOrder.Status == order.StatusSpam && status != order.StatusSpam {
		if err := u.repo.UnblockIP(ctx, oldOrder.IPAddress); err != nil {
			return err
		}
	}

	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *UseCase) MarkAsSpam(ctx context.Context, id int64) error {
	oldOrder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if oldOrder == nil {
		return apperror.New(apperror.CodeOrderNotFound, "Order not found", map[string]any{
			"id": id,
		})
	}

	return u.repo.MarkAsSpam(ctx, id)
}

func (u *UseCase) ExportOrders(ctx context.Context) ([]*order.Order, error) {
	return u.repo.Export(ctx)
}
