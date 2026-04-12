package product

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/infrastructure/gemini"
)

const aiFallbackLimit = 12

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

	if len(candidates) < 10 {
		popular, _, _ := u.repo.List(ctx, product.ListFilter{
			Status: ptr(product.StatusPublished),
			Limit:  40,
			SortBy: "views_count",
		})
		candidates = append(candidates, popular...)
	}

	uniqueCandidates := make([]*product.Project, 0)
	seen := make(map[int64]bool)
	for _, c := range candidates {
		if !seen[c.ID] {
			seen[c.ID] = true
			uniqueCandidates = append(uniqueCandidates, c)
		}
	}

	// Get categories map for names
	cats, _ := u.repo.ListCategories(ctx)
	catMap := make(map[int64]string)
	for _, c := range cats {
		catMap[c.ID] = c.Name
	}
	candidatePool := filterByRequestedCategory(query, uniqueCandidates, catMap)

	// 3. Simplify candidates for Gemini with CATEGORY NAME
	type simpleProj struct {
		ID       int64   `json:"id"`
		Category string  `json:"category"`
		Name     string  `json:"name"`
		Budget   float64 `json:"budget"`
		Tags     string  `json:"tags"`
	}
	simpleProjects := make([]simpleProj, len(candidatePool))
	for i, p := range candidatePool {
		catName := "Прочее"
		if p.ProjectCategoryID != nil {
			catName = catMap[*p.ProjectCategoryID]
		}
		simpleProjects[i] = simpleProj{
			ID:       p.ID,
			Category: catName,
			Name:     p.Name,
			Budget:   p.Budget,
			Tags:     p.AITags,
		}
	}

	productsJSON, _ := json.Marshal(simpleProjects)

	// 4. Gemini selects the BEST 8
	ids, err := u.gemini.SearchProducts(ctx, query, string(productsJSON))
	if err != nil {
		if errors.Is(err, gemini.ErrDisabled) {
			u.logger.Info("Gemini search disabled, using candidates as fallback")
		} else {
			u.logger.Warn("Gemini API error, using candidates as fallback", "error", err)
		}
		return limitProjects(candidatePool, aiFallbackLimit), nil
	}

	u.logger.Info("Gemini returned project IDs", "ids", ids)

	if len(ids) == 0 {
		return []*product.Project{}, nil
	}

	var results []*product.Project
	for _, id := range ids {
		p, err := u.repo.GetByID(ctx, id)
		if err == nil && p != nil && p.Status == product.StatusPublished && matchesRequestedCategory(query, p, catMap) {
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

func limitProjects(projects []*product.Project, limit int) []*product.Project {
	if limit <= 0 || len(projects) <= limit {
		return projects
	}
	return projects[:limit]
}

func filterByRequestedCategory(query string, projects []*product.Project, categories map[int64]string) []*product.Project {
	needle := requestedCategoryNeedle(query)
	if needle == "" {
		return projects
	}

	filtered := make([]*product.Project, 0, len(projects))
	for _, project := range projects {
		if projectMatchesCategory(project, categories, needle) {
			filtered = append(filtered, project)
		}
	}
	if len(filtered) == 0 {
		return projects
	}
	return filtered
}

func matchesRequestedCategory(query string, project *product.Project, categories map[int64]string) bool {
	needle := requestedCategoryNeedle(query)
	return needle == "" || projectMatchesCategory(project, categories, needle)
}

func requestedCategoryNeedle(query string) string {
	query = strings.ToLower(query)
	switch {
	case strings.Contains(query, "кух"):
		return "кух"
	case strings.Contains(query, "шкаф"), strings.Contains(query, "гардероб"):
		return "шкаф"
	default:
		return ""
	}
}

func projectMatchesCategory(project *product.Project, categories map[int64]string, needle string) bool {
	if project == nil || project.ProjectCategoryID == nil {
		return false
	}
	categoryName := strings.ToLower(categories[*project.ProjectCategoryID])
	return strings.Contains(categoryName, needle)
}
