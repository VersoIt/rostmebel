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

func (u *AIUseCase) Search(ctx context.Context, query string) ([]*product.Project, error) {
	u.logger.Info("AI Search request", "query", query)

	// Check cache
	cacheKey := fmt.Sprintf("ai_search:%s", hashQuery(query))
	if cached, err := u.redis.Get(ctx, cacheKey).Result(); err == nil {
		var projects []*product.Project
		if err := json.Unmarshal([]byte(cached), &projects); err == nil {
			u.logger.Info("AI Search cache hit", "query", query)
			return projects, nil
		}
	}

	// 1. Fetch Candidates (RAG Lite approach)
	candidates, err := u.repo.Search(ctx, query, 40)
	if err != nil {
		u.logger.Warn("Initial FTS search failed", "error", err)
	}

	// 2. If FTS didn't find enough, add some popular/recent projects as context
	if len(candidates) < 10 {
		popular, _, _ := u.repo.List(ctx, product.ListFilter{
			Status: ptr(product.StatusPublished),
			Limit:  40,
			SortBy: "views_count",
		})
		candidates = append(candidates, popular...)
	}

	// Remove duplicates
	uniqueCandidates := make([]*product.Project, 0)
	seen := make(map[int64]bool)
	for _, c := range candidates {
		if !seen[c.ID] {
			seen[c.ID] = true
			uniqueCandidates = append(uniqueCandidates, c)
		}
	}

	// 3. Simplify candidates for Gemini
	type simpleProj struct {
		ID     int64   `json:"id"`
		Name   string  `json:"name"`
		Budget float64 `json:"budget"`
		Tags   string  `json:"tags"`
	}
	simpleProjects := make([]simpleProj, len(uniqueCandidates))
	for i, p := range uniqueCandidates {
		simpleProjects[i] = simpleProj{
			ID:     p.ID,
			Name:   p.Name,
			Budget: p.Budget,
			Tags:   p.AITags,
		}
	}

	productsJSON, _ := json.Marshal(simpleProjects)
	
	// 4. Gemini selects the BEST 8
	ids, err := u.gemini.SearchProducts(ctx, query, string(productsJSON))
	if err != nil {
		u.logger.Warn("Gemini API error, using candidates as fallback", "error", err)
		return uniqueCandidates, nil
	}

	u.logger.Info("Gemini returned project IDs", "ids", ids)

	if len(ids) == 0 {
		return []*product.Project{}, nil
	}

	var results []*product.Project
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
