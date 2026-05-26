package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appProduct "github.com/rostmebel/backend/internal/application/product"
	domProduct "github.com/rostmebel/backend/internal/domain/product"
)

type sitemapProductRepo struct {
	projects []*domProduct.Project
}

func (r *sitemapProductRepo) GetByID(context.Context, int64) (*domProduct.Project, error) {
	return nil, nil
}

func (r *sitemapProductRepo) GetBySlug(context.Context, string) (*domProduct.Project, error) {
	return nil, nil
}

func (r *sitemapProductRepo) List(_ context.Context, filter domProduct.ListFilter) ([]*domProduct.Project, int64, error) {
	var projects []*domProduct.Project
	for _, project := range r.projects {
		if filter.Status != nil && project.Status != *filter.Status {
			continue
		}
		projects = append(projects, project)
	}
	return projects, int64(len(projects)), nil
}

func (r *sitemapProductRepo) Create(context.Context, *domProduct.Project) error {
	return nil
}

func (r *sitemapProductRepo) Update(context.Context, *domProduct.Project) error {
	return nil
}

func (r *sitemapProductRepo) Delete(context.Context, int64) error {
	return nil
}

func (r *sitemapProductRepo) ListCategories(context.Context) ([]*domProduct.Category, error) {
	return nil, nil
}

func (r *sitemapProductRepo) GetCategoryBySlug(context.Context, string) (*domProduct.Category, error) {
	return nil, nil
}

func (r *sitemapProductRepo) CreateCategory(context.Context, *domProduct.Category) error {
	return nil
}

func (r *sitemapProductRepo) UpdateCategory(context.Context, *domProduct.Category) error {
	return nil
}

func (r *sitemapProductRepo) DeleteCategory(context.Context, int64) error {
	return nil
}

func (r *sitemapProductRepo) IncrementViews(context.Context, int64) error {
	return nil
}

func (r *sitemapProductRepo) Search(context.Context, string, int) ([]*domProduct.Project, error) {
	return nil, nil
}

func TestSitemapIncludesPublishedProjectSlugs(t *testing.T) {
	updatedAt := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	repo := &sitemapProductRepo{
		projects: []*domProduct.Project{
			{
				ID:        1,
				Slug:      "published-kitchen",
				Status:    domProduct.StatusPublished,
				UpdatedAt: updatedAt,
				Name:      "Published Kitchen",
				Images: []domProduct.Image{
					{URL: "/uploads/kitchen-main.jpg", IsMain: true},
					{URL: "https://cdn.example.com/kitchen-side.jpg"},
				},
			},
			{ID: 2, Slug: "draft-project", Status: domProduct.StatusDraft, UpdatedAt: updatedAt},
		},
	}
	handler := NewProductHandler(appProduct.NewUseCase(repo), nil, "https://example.com/")
	rec := httptest.NewRecorder()

	handler.Sitemap(rec, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<loc>https://example.com/product/published-kitchen</loc>") {
		t.Fatalf("expected published project slug in sitemap, got %s", body)
	}
	if !strings.Contains(body, "<image:image>") {
		t.Fatalf("expected image sitemap entries, got %s", body)
	}
	if !strings.Contains(body, "<image:loc>https://example.com/uploads/kitchen-main.jpg</image:loc>") {
		t.Fatalf("expected relative image URL to be absolutized, got %s", body)
	}
	if !strings.Contains(body, "<image:loc>https://cdn.example.com/kitchen-side.jpg</image:loc>") {
		t.Fatalf("expected absolute image URL in sitemap, got %s", body)
	}
	if strings.Contains(body, "draft-project") {
		t.Fatalf("did not expect draft project in sitemap, got %s", body)
	}
	if !strings.Contains(body, "<lastmod>2026-04-12</lastmod>") {
		t.Fatalf("expected project lastmod, got %s", body)
	}
}

func TestRobotsReferencesConfiguredSitemap(t *testing.T) {
	handler := NewProductHandler(appProduct.NewUseCase(&sitemapProductRepo{}), nil, "https://example.com/")
	rec := httptest.NewRecorder()

	handler.Robots(rec, httptest.NewRequest(http.MethodGet, "/robots.txt", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Disallow: /admin") {
		t.Fatalf("expected admin disallow in robots, got %s", body)
	}
	if !strings.Contains(body, "Sitemap: https://example.com/sitemap.xml") {
		t.Fatalf("expected configured sitemap URL, got %s", body)
	}
}

func TestRobotsUpgradesConfiguredURLToHTTPSForMatchingHost(t *testing.T) {
	handler := NewProductHandler(appProduct.NewUseCase(&sitemapProductRepo{}), nil, "http://rostmebel.shop")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	req.Host = "rostmebel.shop"
	req.Header.Set("X-Forwarded-Host", "rostmebel.shop")
	req.Header.Set("X-Forwarded-Proto", "https")

	handler.Robots(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Sitemap: https://rostmebel.shop/sitemap.xml") {
		t.Fatalf("expected https sitemap URL, got %s", body)
	}
}

func TestRobotsKeepsConfiguredURLForDifferentHost(t *testing.T) {
	handler := NewProductHandler(appProduct.NewUseCase(&sitemapProductRepo{}), nil, "http://rostmebel.shop")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	req.Host = "preview.rostmebel.shop"
	req.Header.Set("X-Forwarded-Host", "preview.rostmebel.shop")
	req.Header.Set("X-Forwarded-Proto", "https")

	handler.Robots(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Sitemap: http://rostmebel.shop/sitemap.xml") {
		t.Fatalf("expected configured sitemap URL for mismatched host, got %s", body)
	}
}
