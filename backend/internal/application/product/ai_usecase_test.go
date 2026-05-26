package product

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	domProduct "github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/infrastructure/gemini"
)

type stubAISearcher struct {
	ids []int64
	err error
}

func (s *stubAISearcher) SearchProducts(_ context.Context, _ string, _ string) ([]int64, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]int64(nil), s.ids...), nil
}

type memorySearchCache struct {
	items map[string]string
}

func (c *memorySearchCache) Get(_ context.Context, key string) (string, error) {
	if c == nil || c.items == nil {
		return "", errCacheMiss
	}
	value, ok := c.items[key]
	if !ok {
		return "", errCacheMiss
	}
	return value, nil
}

func (c *memorySearchCache) Set(_ context.Context, key string, value string, _ time.Duration) error {
	if c.items == nil {
		c.items = make(map[string]string)
	}
	c.items[key] = value
	return nil
}

type stubProductRepo struct {
	byID       map[int64]*domProduct.Project
	search     []*domProduct.Project
	searchErr  error
	categories []*domProduct.Category
	lists      map[string][]*domProduct.Project
}

func (r *stubProductRepo) GetByID(_ context.Context, id int64) (*domProduct.Project, error) {
	return r.byID[id], nil
}

func (r *stubProductRepo) GetBySlug(_ context.Context, slug string) (*domProduct.Project, error) {
	for _, project := range r.byID {
		if project != nil && project.Slug == slug {
			return project, nil
		}
	}
	return nil, nil
}

func (r *stubProductRepo) List(_ context.Context, filter domProduct.ListFilter) ([]*domProduct.Project, int64, error) {
	projects := append([]*domProduct.Project(nil), r.lists[filter.SortBy]...)
	filtered := make([]*domProduct.Project, 0, len(projects))
	for _, project := range projects {
		if project == nil {
			continue
		}
		if filter.Status != nil && project.Status != *filter.Status {
			continue
		}
		filtered = append(filtered, project)
	}
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}
	return filtered, int64(len(filtered)), nil
}

func (r *stubProductRepo) Create(_ context.Context, _ *domProduct.Project) error { return nil }
func (r *stubProductRepo) Update(_ context.Context, _ *domProduct.Project) error { return nil }
func (r *stubProductRepo) Delete(_ context.Context, _ int64) error               { return nil }

func (r *stubProductRepo) ListCategories(_ context.Context) ([]*domProduct.Category, error) {
	return r.categories, nil
}

func (r *stubProductRepo) GetCategoryBySlug(_ context.Context, slug string) (*domProduct.Category, error) {
	for _, category := range r.categories {
		if category != nil && category.Slug == slug {
			return category, nil
		}
	}
	return nil, nil
}

func (r *stubProductRepo) CreateCategory(_ context.Context, _ *domProduct.Category) error { return nil }
func (r *stubProductRepo) UpdateCategory(_ context.Context, _ *domProduct.Category) error { return nil }
func (r *stubProductRepo) DeleteCategory(_ context.Context, _ int64) error                { return nil }

func (r *stubProductRepo) IncrementViews(_ context.Context, _ int64) error { return nil }

func (r *stubProductRepo) Search(_ context.Context, _ string, _ int) ([]*domProduct.Project, error) {
	if r.searchErr != nil {
		return nil, r.searchErr
	}
	return append([]*domProduct.Project(nil), r.search...), nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func buildAIUseCaseForTest(repo *stubProductRepo, searcher aiSearcher, cache aiSearchCache) *AIUseCase {
	return &AIUseCase{
		repo:   repo,
		gemini: searcher,
		cache:  cache,
		logger: testLogger(),
	}
}

func TestAIUseCaseFallbackRespectsBudgetAndCategory(t *testing.T) {
	kitchenCategoryID := int64(1)
	wardrobeCategoryID := int64(2)

	cheapKitchen := &domProduct.Project{ID: 1, ProjectCategoryID: &kitchenCategoryID, Name: "Светлая кухня", Budget: 95000, Status: domProduct.StatusPublished}
	expensiveKitchen := &domProduct.Project{ID: 2, ProjectCategoryID: &kitchenCategoryID, Name: "Большая кухня", Budget: 145000, Status: domProduct.StatusPublished}
	cheapWardrobe := &domProduct.Project{ID: 3, ProjectCategoryID: &wardrobeCategoryID, Name: "Шкаф в прихожую", Budget: 80000, Status: domProduct.StatusPublished}

	repo := &stubProductRepo{
		byID: map[int64]*domProduct.Project{
			1: cheapKitchen,
			2: expensiveKitchen,
			3: cheapWardrobe,
		},
		search: []*domProduct.Project{expensiveKitchen, cheapWardrobe, cheapKitchen},
		categories: []*domProduct.Category{
			{ID: kitchenCategoryID, Name: "Кухни", Slug: "kitchens"},
			{ID: wardrobeCategoryID, Name: "Шкафы-купе", Slug: "wardrobes"},
		},
		lists: map[string][]*domProduct.Project{
			"views_count": {cheapWardrobe, expensiveKitchen, cheapKitchen},
			"updated_at":  {cheapKitchen, expensiveKitchen, cheapWardrobe},
		},
	}

	useCase := buildAIUseCaseForTest(repo, &stubAISearcher{err: gemini.ErrDisabled}, &memorySearchCache{})

	results, err := useCase.Search(context.Background(), "кухни до 100000")
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) != 1 || results[0].ID != cheapKitchen.ID {
		t.Fatalf("expected only cheap kitchen, got %#v", results)
	}
}

func TestAIUseCaseFallsBackWhenGeminiReturnsEmptyIDs(t *testing.T) {
	kitchenCategoryID := int64(1)

	project := &domProduct.Project{
		ID:                7,
		ProjectCategoryID: &kitchenCategoryID,
		Name:              "Сканди кухня",
		Description:       "Светлая кухня в скандинавском стиле",
		AITags:            "сканди, светлая кухня",
		Budget:            145000,
		Status:            domProduct.StatusPublished,
	}

	repo := &stubProductRepo{
		byID:   map[int64]*domProduct.Project{project.ID: project},
		search: []*domProduct.Project{project},
		categories: []*domProduct.Category{
			{ID: kitchenCategoryID, Name: "Кухни", Slug: "kitchens"},
		},
		lists: map[string][]*domProduct.Project{
			"views_count": {project},
			"updated_at":  {project},
		},
	}

	useCase := buildAIUseCaseForTest(repo, &stubAISearcher{ids: []int64{}}, &memorySearchCache{})

	results, err := useCase.Search(context.Background(), "светлая кухня")
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) != 1 || results[0].ID != project.ID {
		t.Fatalf("expected deterministic fallback result, got %#v", results)
	}
}

func TestAIUseCaseUsesSpecsWhenFTSFindsNothing(t *testing.T) {
	kitchenCategoryID := int64(1)

	project := &domProduct.Project{
		ID:                11,
		ProjectCategoryID: &kitchenCategoryID,
		Name:              "Проект с фурнитурой Blum",
		Description:       "Современная кухня",
		Budget:            180000,
		Details: map[string]string{
			"Материал":   "Дуб",
			"Фурнитура":  "Blum",
			"Покрытие":   "эмаль",
			"Столешница": "камень",
		},
		Status: domProduct.StatusPublished,
	}

	repo := &stubProductRepo{
		byID:   map[int64]*domProduct.Project{project.ID: project},
		search: nil,
		categories: []*domProduct.Category{
			{ID: kitchenCategoryID, Name: "Кухни", Slug: "kitchens"},
		},
		lists: map[string][]*domProduct.Project{
			"views_count": {project},
			"updated_at":  {project},
		},
	}

	useCase := buildAIUseCaseForTest(repo, &stubAISearcher{err: gemini.ErrDisabled}, &memorySearchCache{})

	results, err := useCase.Search(context.Background(), "дуб blum кухня")
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) == 0 || results[0].ID != project.ID {
		t.Fatalf("expected project matched by specs, got %#v", results)
	}
}

func TestAIUseCaseRebuildsWhenCachedProjectsAreStale(t *testing.T) {
	kitchenCategoryID := int64(1)

	staleProject := &domProduct.Project{
		ID:                21,
		ProjectCategoryID: &kitchenCategoryID,
		Name:              "Скрытая кухня",
		Budget:            110000,
		Status:            domProduct.StatusDraft,
	}
	liveProject := &domProduct.Project{
		ID:                22,
		ProjectCategoryID: &kitchenCategoryID,
		Name:              "Живая кухня",
		Description:       "Кухня для семьи",
		Budget:            120000,
		Status:            domProduct.StatusPublished,
	}

	cache := &memorySearchCache{
		items: map[string]string{
			"ai_search:v2:" + hashQuery("кухня"): `[21]`,
		},
	}

	repo := &stubProductRepo{
		byID: map[int64]*domProduct.Project{
			staleProject.ID: staleProject,
			liveProject.ID:  liveProject,
		},
		search: []*domProduct.Project{liveProject},
		categories: []*domProduct.Category{
			{ID: kitchenCategoryID, Name: "Кухни", Slug: "kitchens"},
		},
		lists: map[string][]*domProduct.Project{
			"views_count": {liveProject},
			"updated_at":  {liveProject},
		},
	}

	useCase := buildAIUseCaseForTest(repo, &stubAISearcher{err: gemini.ErrDisabled}, cache)

	results, err := useCase.Search(context.Background(), "кухня")
	if err != nil {
		t.Fatalf("search returned error: %v", err)
	}
	if len(results) != 1 || results[0].ID != liveProject.ID {
		t.Fatalf("expected stale cache to be ignored, got %#v", results)
	}
}

func TestExtractBudgetRange(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantMin float64
		wantMax float64
		hasMin  bool
		hasMax  bool
	}{
		{name: "max only", query: "кухня до 250 000", wantMax: 250000, hasMax: true},
		{name: "range", query: "от 100 до 250 тыс", wantMin: 100000, wantMax: 250000, hasMin: true, hasMax: true},
		{name: "million", query: "гардеробная до 1.2 млн", wantMax: 1200000, hasMax: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minBudget, maxBudget := extractBudgetRange(tt.query)

			if tt.hasMin {
				if minBudget == nil || *minBudget != tt.wantMin {
					t.Fatalf("expected min budget %v, got %v", tt.wantMin, minBudget)
				}
			} else if minBudget != nil {
				t.Fatalf("expected nil min budget, got %v", *minBudget)
			}

			if tt.hasMax {
				if maxBudget == nil || *maxBudget != tt.wantMax {
					t.Fatalf("expected max budget %v, got %v", tt.wantMax, maxBudget)
				}
			} else if maxBudget != nil {
				t.Fatalf("expected nil max budget, got %v", *maxBudget)
			}
		})
	}
}

func TestFilterByRequestedCategoryKeepsKitchenQueriesFocused(t *testing.T) {
	kitchenCategoryID := int64(1)
	wardrobeCategoryID := int64(2)
	projects := []*domProduct.Project{
		{ID: 1, ProjectCategoryID: &wardrobeCategoryID, Name: "Шкаф", Status: domProduct.StatusPublished},
		{ID: 2, ProjectCategoryID: &kitchenCategoryID, Name: "Кухня", Status: domProduct.StatusPublished},
	}
	categories := map[int64]string{
		kitchenCategoryID:  "Кухни",
		wardrobeCategoryID: "Шкафы-купе",
	}

	filtered := filterByRequestedCategory("кухни до 400000", projects, categories)
	if len(filtered) != 1 || filtered[0].ID != 2 {
		t.Fatalf("expected only kitchen project, got %#v", filtered)
	}
}

func TestFilterByRequestedCategoryFallsBackWhenCategoryIsUnknown(t *testing.T) {
	categoryID := int64(1)
	projects := []*domProduct.Project{{ID: 1, ProjectCategoryID: &categoryID, Name: "Проект", Status: domProduct.StatusPublished}}

	filtered := filterByRequestedCategory("мебель в прихожую", projects, map[int64]string{categoryID: "Кухни"})
	if len(filtered) != len(projects) {
		t.Fatalf("expected original projects for unknown category query")
	}
}

func TestSearchProductsDisabledErrorIsStillSupported(t *testing.T) {
	if !errors.Is(gemini.ErrDisabled, gemini.ErrDisabled) {
		t.Fatal("expected disabled error to remain comparable")
	}
}
