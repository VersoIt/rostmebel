package review

import (
	"context"
	"fmt"
	"strings"

	"github.com/rostmebel/backend/internal/domain/review"
	"github.com/rostmebel/backend/internal/domain/order"
)

type UseCase struct {
	repo      review.Repository
	orderRepo order.Repository
}

func NewUseCase(repo review.Repository, orderRepo order.Repository) *UseCase {
	return &UseCase{repo: repo, orderRepo: orderRepo}
}

func (u *UseCase) CreateReview(ctx context.Context, clientPhone string, rev *review.Review) error {
	// 1. Normalize phone for search
	cleanPhone := normalizePhone(clientPhone)
	if cleanPhone == "" {
		return fmt.Errorf("некорректный номер телефона")
	}

	// 2. Find order by phone with status 'done'
	orders, _, _, err := u.orderRepo.List(ctx, order.ListFilter{
		Limit: 100, // Enough to find recent
	})
	if err != nil {
		return err
	}

	var validOrder *order.Order
	for _, o := range orders {
		// Basic check: suffix match or exact match after normalization
		if normalizePhone(o.ClientPhone) == cleanPhone && o.Status == order.StatusDone {
			validOrder = o
			break
		}
	}

	if validOrder == nil {
		return fmt.Errorf("вы не можете оставить отзыв: заказ на данный номер не найден или еще не завершен")
	}

	// 3. Link review to order and set pending
	rev.OrderID = validOrder.ID
	rev.Status = review.StatusPending
	
	// If project_id not provided, try to take from order
	if rev.ProjectID == nil && validOrder.ProjectID != nil {
		rev.ProjectID = validOrder.ProjectID
	}

	return u.repo.Create(ctx, rev)
}

func (u *UseCase) ListReviews(ctx context.Context, f review.ListFilter) ([]*review.Review, int64, int64, error) {
	return u.repo.List(ctx, f)
}

func (u *UseCase) GetByProject(ctx context.Context, projectID int64) ([]*review.Review, error) {
	return u.repo.GetByProjectID(ctx, projectID)
}

func (u *UseCase) ModerateReview(ctx context.Context, id int64, approved bool) error {
	status := review.StatusRejected
	if approved {
		status = review.StatusApproved
	}
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *UseCase) DeleteReview(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func normalizePhone(p string) string {
	var b strings.Builder
	for _, r := range p {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	s := b.String()
	if len(s) > 10 {
		return s[len(s)-10:] // Get last 10 digits
	}
	return s
}
