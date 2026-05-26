package product

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/redis/go-redis/v9"
	"github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/infrastructure/gemini"
	"golang.org/x/sync/errgroup"
)

const (
	aiFallbackLimit         = 12
	aiGeminiCandidateLimit  = 24
	aiFTSCandidateLimit     = 40
	aiSupplementalListLimit = 80
	aiSearchCacheTTL        = 5 * time.Minute
	aiSearchCacheVersion    = "v2"
)

var (
	errCacheMiss = errors.New("ai search cache miss")

	budgetRangePattern = regexp.MustCompile(`(?i)(?:от\s+)?([0-9][0-9\s.,]*\s*(?:к|k|тыс\.?|тысяч(?:и)?|млн|м)?)\s*(?:-|–|—|до)\s*([0-9][0-9\s.,]*\s*(?:к|k|тыс\.?|тысяч(?:и)?|млн|м)?)`)
	budgetMaxPattern   = regexp.MustCompile(`(?i)(?:до|не\s+дороже|не\s+больше)\s+([0-9][0-9\s.,]*\s*(?:к|k|тыс\.?|тысяч(?:и)?|млн|м)?)`)
	budgetMinPattern   = regexp.MustCompile(`(?i)(?:от|минимум|не\s+меньше)\s+([0-9][0-9\s.,]*\s*(?:к|k|тыс\.?|тысяч(?:и)?|млн|м)?)`)
)

var budgetNoiseTokens = map[string]struct{}{
	"до": {}, "от": {}, "к": {}, "k": {}, "м": {}, "тыс": {}, "тыс.": {}, "тысяч": {}, "тысячи": {},
	"млн": {}, "р": {}, "руб": {}, "рублей": {}, "рубля": {},
}

type aiSearcher interface {
	SearchProducts(ctx context.Context, userQuery string, projectsJSON string) ([]int64, error)
}

type aiSearchCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type redisSearchCache struct {
	client *redis.Client
}

func (c redisSearchCache) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", errCacheMiss
	}

	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errCacheMiss
	}

	return value, err
}

func (c redisSearchCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

type AIUseCase struct {
	repo   product.Repository
	gemini aiSearcher
	cache  aiSearchCache
	logger *slog.Logger
}

type searchConstraints struct {
	MinBudget  *float64
	MaxBudget  *float64
	Categories map[int64]struct{}
}

type simpleProj struct {
	ID          int64   `json:"id"`
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Budget      float64 `json:"budget"`
	Description string  `json:"description"`
	Tags        string  `json:"tags"`
	Specs       string  `json:"specs"`
}

func NewAIUseCase(repo product.Repository, geminiClient aiSearcher, redisClient *redis.Client, logger *slog.Logger) *AIUseCase {
	if logger == nil {
		logger = slog.Default()
	}

	return &AIUseCase{
		repo:   repo,
		gemini: geminiClient,
		cache:  redisSearchCache{client: redisClient},
		logger: logger,
	}
}

func (u *AIUseCase) Search(ctx context.Context, rawQuery string) ([]*product.Project, error) {
	query := normalizeSearchQuery(rawQuery)
	if query == "" {
		return []*product.Project{}, nil
	}

	u.logger.Info("AI Search request", "query", query)

	categories, _ := u.repo.ListCategories(ctx)
	categoryMap := categoryNameMap(categories)
	constraints := buildSearchConstraintsFromCategories(query, categories)
	cacheKey := fmt.Sprintf("ai_search:%s:%s", aiSearchCacheVersion, hashQuery(query))

	if cachedIDs, ok := u.readCachedIDs(ctx, cacheKey); ok {
		if len(cachedIDs) == 0 {
			u.logger.Info("AI Search cache hit with empty result", "query", query)
			return []*product.Project{}, nil
		}

		cachedResults := u.loadProjectsByID(ctx, cachedIDs, constraints)
		if len(cachedResults) > 0 {
			u.logger.Info("AI Search cache hit", "query", query, "count", len(cachedResults))
			return cachedResults, nil
		}

		u.logger.Warn("AI Search cache entry became stale, rebuilding", "query", query)
	}

	ftsQuery := buildFTSQuery(query)
	candidates := u.collectCandidates(ctx, ftsQuery)
	candidatePool := filterProjects(candidates, constraints)
	fallbackResults := limitProjects(candidatePool, aiFallbackLimit)

	if len(candidatePool) == 0 {
		u.writeCachedIDs(ctx, cacheKey, nil)
		return []*product.Project{}, nil
	}

	if u.gemini == nil {
		u.logger.Info("Gemini client unavailable, using deterministic fallback")
		u.writeCachedIDs(ctx, cacheKey, projectIDs(fallbackResults))
		return fallbackResults, nil
	}

	productsJSON, err := json.Marshal(buildGeminiCandidates(candidatePool, categoryMap, aiGeminiCandidateLimit))
	if err != nil {
		u.logger.Warn("Failed to marshal Gemini candidates, using deterministic fallback", "error", err)
		u.writeCachedIDs(ctx, cacheKey, projectIDs(fallbackResults))
		return fallbackResults, nil
	}

	ids, err := u.gemini.SearchProducts(ctx, query, string(productsJSON))
	if err != nil {
		if errors.Is(err, gemini.ErrDisabled) {
			u.logger.Info("Gemini search disabled, using deterministic fallback")
		} else {
			u.logger.Warn("Gemini API error, using deterministic fallback", "error", err)
		}

		u.writeCachedIDs(ctx, cacheKey, projectIDs(fallbackResults))
		return fallbackResults, nil
	}

	primaryResults := u.loadProjectsByID(ctx, ids, constraints)
	finalResults := mergeProjectLists(primaryResults, fallbackResults, aiFallbackLimit)
	if len(finalResults) == 0 {
		u.logger.Info("Gemini returned no usable projects, using fallback", "query", query)
		finalResults = fallbackResults
	}

	u.writeCachedIDs(ctx, cacheKey, projectIDs(finalResults))
	return finalResults, nil
}

func (u *AIUseCase) readCachedIDs(ctx context.Context, cacheKey string) ([]int64, bool) {
	if u.cache == nil {
		return nil, false
	}

	cached, err := u.cache.Get(ctx, cacheKey)
	if err != nil {
		if !errors.Is(err, errCacheMiss) {
			u.logger.Warn("AI Search cache read failed", "key", cacheKey, "error", err)
		}
		return nil, false
	}

	var ids []int64
	if err := json.Unmarshal([]byte(cached), &ids); err != nil {
		u.logger.Warn("AI Search cache payload is invalid", "key", cacheKey, "error", err)
		return nil, false
	}

	return ids, true
}

func (u *AIUseCase) writeCachedIDs(ctx context.Context, cacheKey string, ids []int64) {
	if u.cache == nil {
		return
	}

	payload, err := json.Marshal(ids)
	if err != nil {
		u.logger.Warn("AI Search cache marshal failed", "key", cacheKey, "error", err)
		return
	}

	if err := u.cache.Set(ctx, cacheKey, string(payload), aiSearchCacheTTL); err != nil {
		u.logger.Warn("AI Search cache write failed", "key", cacheKey, "error", err)
	}
}

func (u *AIUseCase) collectCandidates(ctx context.Context, query string) []*product.Project {
	unique := make([]*product.Project, 0, aiFTSCandidateLimit+(2*aiSupplementalListLimit))
	seen := make(map[int64]struct{}, cap(unique))
	var mu sync.Mutex

	addProjects := func(projects []*product.Project) {
		for _, project := range projects {
			if project == nil || project.Status != product.StatusPublished {
				continue
			}

			mu.Lock()
			if _, ok := seen[project.ID]; ok {
				mu.Unlock()
				continue
			}
			seen[project.ID] = struct{}{}
			unique = append(unique, project)
			mu.Unlock()
		}
	}

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		projects, err := u.repo.Search(groupCtx, query, aiFTSCandidateLimit)
		if err != nil {
			u.logger.Warn("Initial FTS search failed", "query", query, "error", err)
			return nil
		}
		addProjects(projects)
		return nil
	})

	group.Go(func() error {
		projects, err := u.listPublished(groupCtx, "views_count")
		if err != nil {
			u.logger.Warn("Supplemental project list failed", "sort_by", "views_count", "error", err)
			return nil
		}
		addProjects(projects)
		return nil
	})

	group.Go(func() error {
		projects, err := u.listPublished(groupCtx, "updated_at")
		if err != nil {
			u.logger.Warn("Supplemental project list failed", "sort_by", "updated_at", "error", err)
			return nil
		}
		addProjects(projects)
		return nil
	})

	_ = group.Wait()
	return unique
}

func (u *AIUseCase) listPublished(ctx context.Context, sortBy string) ([]*product.Project, error) {
	status := product.StatusPublished
	projects, _, err := u.repo.List(ctx, product.ListFilter{
		Status:    &status,
		Limit:     aiSupplementalListLimit,
		SortBy:    sortBy,
		SortOrder: "DESC",
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (u *AIUseCase) loadProjectsByID(ctx context.Context, ids []int64, constraints searchConstraints) []*product.Project {
	results := make([]*product.Project, 0, len(ids))
	for _, id := range ids {
		project, err := u.repo.GetByID(ctx, id)
		if err != nil {
			u.logger.Warn("Failed to load AI Search project", "project_id", id, "error", err)
			continue
		}
		if !projectPassesConstraints(project, constraints) {
			continue
		}
		results = append(results, project)
	}

	return results
}

func buildGeminiCandidates(projects []*product.Project, categories map[int64]string, limit int) []simpleProj {
	if limit <= 0 || len(projects) < limit {
		limit = len(projects)
	}

	result := make([]simpleProj, 0, limit)
	for _, project := range projects[:limit] {
		if project == nil {
			continue
		}

		categoryName := "Прочее"
		if project.ProjectCategoryID != nil {
			if name := categories[*project.ProjectCategoryID]; name != "" {
				categoryName = name
			}
		}

		result = append(result, simpleProj{
			ID:          project.ID,
			Category:    categoryName,
			Name:        project.Name,
			Budget:      project.Budget,
			Description: compactProjectText(project.Description),
			Tags:        compactProjectText(project.AITags),
			Specs:       flattenProjectDetails(project.Details),
		})
	}

	return result
}

func mergeProjectLists(primary []*product.Project, secondary []*product.Project, limit int) []*product.Project {
	if limit <= 0 {
		return nil
	}

	merged := make([]*product.Project, 0, limit)
	seen := make(map[int64]struct{}, limit)

	appendUnique := func(projects []*product.Project) {
		for _, project := range projects {
			if project == nil {
				continue
			}
			if _, ok := seen[project.ID]; ok {
				continue
			}
			seen[project.ID] = struct{}{}
			merged = append(merged, project)
			if len(merged) >= limit {
				return
			}
		}
	}

	appendUnique(primary)
	if len(merged) < limit {
		appendUnique(secondary)
	}

	return merged
}

func limitProjects(projects []*product.Project, limit int) []*product.Project {
	if limit <= 0 || len(projects) <= limit {
		return projects
	}

	return projects[:limit]
}

func filterProjects(projects []*product.Project, constraints searchConstraints) []*product.Project {
	filtered := make([]*product.Project, 0, len(projects))
	for _, project := range projects {
		if projectPassesConstraints(project, constraints) {
			filtered = append(filtered, project)
		}
	}
	return filtered
}

func projectPassesConstraints(project *product.Project, constraints searchConstraints) bool {
	if project == nil || project.Status != product.StatusPublished {
		return false
	}

	if len(constraints.Categories) > 0 {
		if project.ProjectCategoryID == nil {
			return false
		}
		if _, ok := constraints.Categories[*project.ProjectCategoryID]; !ok {
			return false
		}
	}

	if constraints.MinBudget != nil && project.Budget < *constraints.MinBudget {
		return false
	}
	if constraints.MaxBudget != nil && project.Budget > *constraints.MaxBudget {
		return false
	}

	return true
}

func buildSearchConstraintsFromCategories(query string, categories []*product.Category) searchConstraints {
	minBudget, maxBudget := extractBudgetRange(query)
	return searchConstraints{
		MinBudget:  minBudget,
		MaxBudget:  maxBudget,
		Categories: requestedCategoryIDs(query, buildCategoryLookupFromCategories(categories)),
	}
}

func buildSearchConstraints(query string, categoryNames map[int64]string) searchConstraints {
	minBudget, maxBudget := extractBudgetRange(query)
	return searchConstraints{
		MinBudget:  minBudget,
		MaxBudget:  maxBudget,
		Categories: requestedCategoryIDs(query, buildCategoryLookupFromNames(categoryNames)),
	}
}

func buildFTSQuery(query string) string {
	query = normalizeSearchQuery(query)
	if query == "" {
		return ""
	}

	cleaned := budgetRangePattern.ReplaceAllString(query, " ")
	cleaned = budgetMaxPattern.ReplaceAllString(cleaned, " ")
	cleaned = budgetMinPattern.ReplaceAllString(cleaned, " ")
	cleaned = normalizeSearchQuery(cleaned)
	if cleaned == "" {
		return query
	}

	tokens := splitQueryTokens(cleaned)
	if len(tokens) == 0 {
		return query
	}

	filtered := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := budgetNoiseTokens[token]; ok {
			continue
		}
		if _, err := strconv.ParseFloat(strings.ReplaceAll(token, ",", "."), 64); err == nil {
			continue
		}
		filtered = append(filtered, token)
	}

	if len(filtered) == 0 {
		return query
	}

	return strings.Join(filtered, " ")
}

func extractBudgetRange(query string) (*float64, *float64) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil, nil
	}

	if matches := budgetRangePattern.FindStringSubmatch(query); len(matches) == 3 {
		leftRaw, rightRaw := normalizeBudgetRangeSides(matches[1], matches[2])
		left, leftOK := parseBudgetValue(leftRaw)
		right, rightOK := parseBudgetValue(rightRaw)
		if leftOK && rightOK {
			if left > right {
				left, right = right, left
			}
			return ptr(left), ptr(right)
		}
	}

	var minBudget *float64
	var maxBudget *float64

	if matches := budgetMaxPattern.FindStringSubmatch(query); len(matches) == 2 {
		if value, ok := parseBudgetValue(matches[1]); ok {
			maxBudget = ptr(value)
		}
	}

	if matches := budgetMinPattern.FindStringSubmatch(query); len(matches) == 2 {
		if value, ok := parseBudgetValue(matches[1]); ok {
			minBudget = ptr(value)
		}
	}

	return minBudget, maxBudget
}

func normalizeBudgetRangeSides(left string, right string) (string, string) {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	switch {
	case !budgetValueHasUnit(left) && budgetValueHasUnit(right):
		left += " " + budgetValueUnit(right)
	case budgetValueHasUnit(left) && !budgetValueHasUnit(right):
		right += " " + budgetValueUnit(left)
	}

	return left, right
}

func budgetValueHasUnit(value string) bool {
	return budgetValueUnit(value) != ""
}

func budgetValueUnit(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch {
	case strings.Contains(value, "млн"):
		return "млн"
	case strings.Contains(value, "тыс"):
		return "тыс"
	case strings.HasSuffix(value, "к"), strings.HasSuffix(value, "k"):
		return "к"
	case strings.HasSuffix(value, "м"):
		return "м"
	default:
		return ""
	}
}

func parseBudgetValue(raw string) (float64, bool) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return 0, false
	}

	multiplier := 1.0
	switch {
	case strings.Contains(value, "млн"), strings.HasSuffix(value, "м"):
		multiplier = 1_000_000
	case strings.Contains(value, "тыс"), strings.HasSuffix(value, "к"), strings.HasSuffix(value, "k"):
		multiplier = 1_000
	}

	replacer := strings.NewReplacer("тыс.", "", "тыс", "", "тысяч", "", "млн", "", "к", "", "k", "", "м", "")
	value = replacer.Replace(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, ",", ".")

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}

	return number * multiplier, true
}

func requestedCategoryIDs(query string, lookup map[string]map[int64]struct{}) map[int64]struct{} {
	query = normalizeSearchQuery(query)
	if query == "" || len(lookup) == 0 {
		return nil
	}

	requested := make(map[int64]struct{})
	for _, segment := range querySegments(query) {
		ids, ok := lookup[segment]
		if !ok {
			continue
		}
		for id := range ids {
			requested[id] = struct{}{}
		}
	}

	if len(requested) == 0 {
		return nil
	}

	return requested
}

func buildCategoryLookupFromCategories(categories []*product.Category) map[string]map[int64]struct{} {
	lookup := make(map[string]map[int64]struct{}, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		appendCategoryLookup(lookup, category.ID, category.Name)
		appendCategoryLookup(lookup, category.ID, category.Slug)
	}
	return lookup
}

func buildCategoryLookupFromNames(categoryNames map[int64]string) map[string]map[int64]struct{} {
	lookup := make(map[string]map[int64]struct{}, len(categoryNames))
	for id, name := range categoryNames {
		appendCategoryLookup(lookup, id, name)
	}
	return lookup
}

func appendCategoryLookup(lookup map[string]map[int64]struct{}, categoryID int64, value string) {
	for _, key := range categoryKeys(value) {
		ids, ok := lookup[key]
		if !ok {
			ids = make(map[int64]struct{}, 1)
			lookup[key] = ids
		}
		ids[categoryID] = struct{}{}
	}
}

func categoryKeys(value string) []string {
	normalized := normalizeSearchQuery(value)
	if normalized == "" {
		return nil
	}

	keys := make(map[string]struct{})
	keys[normalized] = struct{}{}

	for _, token := range splitQueryTokens(normalized) {
		if token == "" {
			continue
		}
		keys[token] = struct{}{}
	}

	result := make([]string, 0, len(keys))
	for key := range keys {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func querySegments(query string) []string {
	query = normalizeSearchQuery(query)
	if query == "" {
		return nil
	}

	tokens := splitQueryTokens(query)
	if len(tokens) == 0 {
		return []string{query}
	}

	seen := make(map[string]struct{}, len(tokens)*2)
	segments := make([]string, 0, len(tokens)*2)

	appendSegment := func(segment string) {
		if segment == "" {
			return
		}
		if _, ok := seen[segment]; ok {
			return
		}
		seen[segment] = struct{}{}
		segments = append(segments, segment)
	}

	appendSegment(query)

	for start := 0; start < len(tokens); start++ {
		for end := start + 1; end <= len(tokens); end++ {
			appendSegment(strings.Join(tokens[start:end], " "))
		}
	}

	return segments
}

func splitQueryTokens(query string) []string {
	return strings.FieldsFunc(query, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func projectIDs(projects []*product.Project) []int64 {
	ids := make([]int64, 0, len(projects))
	for _, project := range projects {
		if project == nil {
			continue
		}
		ids = append(ids, project.ID)
	}
	return ids
}

func flattenProjectDetails(details map[string]string) string {
	if len(details) == 0 {
		return ""
	}

	keys := make([]string, 0, len(details))
	for key := range details {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(details[key])
		if value == "" {
			continue
		}
		parts = append(parts, key+": "+value)
	}

	return strings.Join(parts, "; ")
}

func compactProjectText(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) <= 320 {
		return value
	}
	return value[:320]
}

func categoryNameMap(categories []*product.Category) map[int64]string {
	result := make(map[int64]string, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		result[category.ID] = category.Name
	}
	return result
}

func normalizeSearchQuery(query string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(query))), " ")
}

func hashQuery(query string) string {
	h := sha256.New()
	h.Write([]byte(query))
	return hex.EncodeToString(h.Sum(nil))
}

func ptr[T any](v T) *T {
	return &v
}

func filterByRequestedCategory(query string, projects []*product.Project, categories map[int64]string) []*product.Project {
	constraints := buildSearchConstraints(query, categories)
	if len(constraints.Categories) == 0 {
		return projects
	}

	filtered := filterProjects(projects, searchConstraints{Categories: constraints.Categories})
	if len(filtered) == 0 {
		return projects
	}

	return filtered
}

func matchesRequestedCategory(query string, project *product.Project, categories map[int64]string) bool {
	constraints := buildSearchConstraints(query, categories)
	if len(constraints.Categories) == 0 {
		return true
	}

	return projectPassesConstraints(project, searchConstraints{Categories: constraints.Categories})
}
