package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/product"
	"github.com/rostmebel/backend/internal/domain/apperror"
	domProduct "github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/interfaces/dto"
	"github.com/xuri/excelize/v2"
)

type ProductHandler struct {
	useCase       *product.UseCase
	aiUseCase     *product.AIUseCase
	publicSiteURL string
}

func NewProductHandler(useCase *product.UseCase, aiUseCase *product.AIUseCase, publicSiteURL string) *ProductHandler {
	return &ProductHandler{
		useCase:       useCase,
		aiUseCase:     aiUseCase,
		publicSiteURL: strings.TrimRight(publicSiteURL, "/"),
	}
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func (h *ProductHandler) Robots(w http.ResponseWriter, r *http.Request) {
	baseURL := h.siteURL(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /admin/\n\nSitemap: %s/sitemap.xml\n", baseURL)
}

func (h *ProductHandler) Sitemap(w http.ResponseWriter, r *http.Request) {
	status := domProduct.StatusPublished
	projects, _, err := h.useCase.ListProjects(r.Context(), domProduct.ListFilter{
		Status:    &status,
		Limit:     10000,
		SortBy:    "updated_at",
		SortOrder: "DESC",
	})
	if err != nil {
		respondWithError(w, err)
		return
	}

	baseURL := h.siteURL(r)
	today := time.Now().UTC().Format("2006-01-02")
	urls := []sitemapURL{
		{Loc: baseURL + "/", LastMod: today, ChangeFreq: "weekly", Priority: "1.0"},
		{Loc: baseURL + "/catalog", LastMod: today, ChangeFreq: "weekly", Priority: "0.8"},
		{Loc: baseURL + "/contact", LastMod: today, ChangeFreq: "monthly", Priority: "0.6"},
	}

	for _, project := range projects {
		path := fmt.Sprintf("/product/%d", project.ID)
		if project.Slug != "" {
			path = "/product/" + project.Slug
		}

		lastMod := project.UpdatedAt.UTC().Format("2006-01-02")
		if project.UpdatedAt.IsZero() {
			lastMod = today
		}

		urls = append(urls, sitemapURL{
			Loc:        baseURL + path,
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "0.7",
		})
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	})
}

func (h *ProductHandler) siteURL(r *http.Request) string {
	if h.publicSiteURL != "" {
		return h.publicSiteURL
	}

	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "https"
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return strings.TrimRight(scheme+"://"+host, "/")
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	filter := domProduct.ListFilter{
		Limit:  20,
		Offset: 0,
	}

	if catID := r.URL.Query().Get("project_category_id"); catID != "" {
		id, err := strconv.ParseInt(catID, 10, 64)
		if err != nil {
			respondWithError(w, invalidQuery("project_category_id", catID))
			return
		}
		filter.ProjectCategoryID = &id
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := domProduct.ProjectStatus(status)
		filter.Status = &s
	}
	if minP := r.URL.Query().Get("min_price"); minP != "" {
		p, err := strconv.ParseFloat(minP, 64)
		if err != nil {
			respondWithError(w, invalidQuery("min_price", minP))
			return
		}
		filter.MinBudget = &p
	}
	if maxP := r.URL.Query().Get("max_price"); maxP != "" {
		p, err := strconv.ParseFloat(maxP, 64)
		if err != nil {
			respondWithError(w, invalidQuery("max_price", maxP))
			return
		}
		filter.MaxBudget = &p
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 1 || l > 10000 {
			respondWithError(w, invalidQuery("limit", limit))
			return
		}
		filter.Limit = l
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil || o < 0 {
			respondWithError(w, invalidQuery("offset", offset))
			return
		}
		filter.Offset = o
	}
	filter.Search = r.URL.Query().Get("search")
	filter.SortBy = r.URL.Query().Get("sort_by")
	filter.SortOrder = r.URL.Query().Get("sort_order")

	projects, total, err := h.useCase.ListProjects(r.Context(), filter)
	if err != nil {
		respondWithError(w, err)
		return
	}

	res := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = dto.FromProject(p)
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": res,
		"total": total,
	})
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Try slug
		p, err := h.useCase.GetProjectBySlug(r.Context(), idStr)
		if err != nil {
			respondWithError(w, err)
			return
		}
		if p == nil {
			respondWithError(w, projectNotFound(idStr))
			return
		}
		h.useCase.IncrementViews(r.Context(), p.ID)
		respondWithJSON(w, http.StatusOK, dto.FromProject(p))
		return
	}

	p, err := h.useCase.GetProject(r.Context(), id)
	if err != nil {
		respondWithError(w, err)
		return
	}
	if p == nil {
		respondWithError(w, projectNotFound(idStr))
		return
	}
	h.useCase.IncrementViews(r.Context(), p.ID)
	respondWithJSON(w, http.StatusOK, dto.FromProject(p))
}

func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.useCase.ListCategories(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	res := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = dto.FromCategory(c)
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (h *ProductHandler) AISearch(w http.ResponseWriter, r *http.Request) {
	var req dto.AISearchRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, err)
		return
	}

	projects, err := h.aiUseCase.Search(r.Context(), req.Query)
	if err != nil {
		respondWithError(w, err)
		return
	}

	res := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = dto.FromProject(p)
	}

	respondWithJSON(w, http.StatusOK, res)
}
func (h *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, apperror.New(apperror.CodeUploadFileMissing, "Image file is required", nil))
		return
	}
	defer file.Close()

	if handler.Size > 10<<20 {
		respondWithError(w, apperror.New(apperror.CodeUploadFileTooLarge, "Image file is too large", map[string]any{
			"max_bytes": 10 << 20,
		}))
		return
	}

	contentType, ext, err := detectImageType(file)
	if err != nil {
		respondWithError(w, err)
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondWithError(w, err)
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"url":          "/uploads/" + filename,
		"content_type": contentType,
	})
}
func (h *ProductHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
	projects, _, err := h.useCase.ListProjects(r.Context(), domProduct.ListFilter{Limit: 10000})
	if err != nil {
		respondWithError(w, err)
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Проекты"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Название", "Slug", "Бюджет", "Статус", "Просмотры", "Заявки"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, head)
	}

	for i, p := range projects {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Slug)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Budget)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), string(p.Status))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.ViewsCount)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.OrdersCount)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=projects.xlsx")
	if err := f.Write(w); err != nil {
		respondWithError(w, err)
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, err)
		return
	}

	p := &domProduct.Project{
		ProjectCategoryID: req.ProjectCategoryID,
		Name:              req.Name,
		Slug:              req.Slug,
		Description:       req.Description,
		Budget:            req.Budget,
		BudgetOld:         req.BudgetOld,
		Images:            req.Images,
		Details:           req.Details,
		AITags:            req.AITags,
		Status:            req.Status,
	}

	if p.Status == "" {
		p.Status = domProduct.StatusDraft
	}

	if err := h.useCase.CreateProject(r.Context(), p); err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromProject(p))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, apperror.New(apperror.CodeInvalidID, "Invalid project id", map[string]any{"id": idStr}))
		return
	}

	var req dto.CreateProjectRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, err)
		return
	}

	p, err := h.useCase.GetProject(r.Context(), id)
	if err != nil {
		respondWithError(w, err)
		return
	}
	if p == nil {
		respondWithError(w, projectNotFound(idStr))
		return
	}

	p.ProjectCategoryID = req.ProjectCategoryID
	p.Name = req.Name
	p.Slug = req.Slug
	p.Description = req.Description
	p.Budget = req.Budget
	p.BudgetOld = req.BudgetOld
	p.Images = req.Images
	p.Details = req.Details
	p.AITags = req.AITags
	p.Status = req.Status

	if err := h.useCase.UpdateProject(r.Context(), p); err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, dto.FromProject(p))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, apperror.New(apperror.CodeInvalidID, "Invalid project id", map[string]any{"id": idStr}))
		return
	}

	if err := h.useCase.DeleteProject(r.Context(), id); err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func invalidQuery(field, value string) error {
	return apperror.New(apperror.CodeInvalidQuery, "Invalid query parameter", map[string]any{
		"field": field,
		"value": value,
	})
}

func projectNotFound(id string) error {
	return apperror.New(apperror.CodeProjectNotFound, "Project not found", map[string]any{
		"id": id,
	})
}

func detectImageType(file io.ReadSeeker) (string, string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", "", err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", err
	}

	contentType := http.DetectContentType(buffer[:n])
	allowed := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
		"image/gif":  ".gif",
	}
	ext, ok := allowed[contentType]
	if !ok {
		return "", "", apperror.New(apperror.CodeUploadInvalidType, "Unsupported image type", map[string]any{
			"content_type": contentType,
			"allowed":      []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
		})
	}
	return contentType, ext, nil
}
