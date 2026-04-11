package product

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/infrastructure/gemini"
	"log/slog"
)

type AIUseCase struct {
	repo   product.Repository
	gemini *gemini.Client
	redis  *redis.Client
	logger *slog.Logger
}

func NewAIUseCase(repo product.Repository, gemini *gemini.Client, redis *redis.Client, logger *slog.Logger) *AIUseCase {
	return &AIUseCase{
		repo:   repo,
		gemini: gemini,
		redis:  redis,
		logger: logger,
	}
}

func (u *AIUseCase) Search(ctx context.Context, query string) ([]*product.Product, error) {
	u.logger.Info("AI Search request", "query", query)

	// Check cache
	cacheKey := fmt.Sprintf("ai_search:%s", hashQuery(query))
	if cached, err := u.redis.Get(ctx, cacheKey).Result(); err == nil {
		var products []*product.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			u.logger.Info("AI Search cache hit", "query", query)
			return products, nil
		}
	}

	// Get all published products
	allProducts, _, err := u.repo.List(ctx, product.ListFilter{
		Status: ptr(product.StatusPublished),
		Limit:  100,
	})
	if err != nil {
		u.logger.Error("AI Search failed to list products", "error", err)
		return nil, err
	}

	// Simplify products for Gemini to reduce tokens and improve precision
	type simpleProd struct {
		ID    int64   `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Tags  string  `json:"tags"`
	}
	simpleProducts := make([]simpleProd, len(allProducts))
	for i, p := range allProducts {
		simpleProducts[i] = simpleProd{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Tags:  p.AITags,
		}
	}

	productsJSON, _ := json.Marshal(simpleProducts)
	
	ids, err := u.gemini.SearchProducts(ctx, query, string(productsJSON))
	if err != nil {
		u.logger.Warn("Gemini API error, falling back to full-text search", "error", err)
		return u.repo.Search(ctx, query, 8)
	}

	u.logger.Info("Gemini returned product IDs", "ids", ids)

	if len(ids) == 0 {
		return []*product.Product{}, nil
	}

	var results []*product.Product
	for _, id := range ids {
		p, err := u.repo.GetByID(ctx, id)
		if err == nil && p != nil && p.Status == product.StatusPublished {
			results = append(results, p)
		}
	}

	if len(results) > 0 {
		data, _ := json.Marshal(results)
		u.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return results, nil
}

func hashQuery(query string) string {
	h := sha256.New()
	h.Write([]byte(query))
	return hex.EncodeToString(h.Sum(nil))
}

func ptr[T any](v T) *T {
	return &v
}
