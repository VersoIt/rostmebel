package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/product"
	domProduct "github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/interfaces/dto"
	"github.com/xuri/excelize/v2"
)

type ProductHandler struct {
	useCase   *product.UseCase
	aiUseCase *product.AIUseCase
}

func NewProductHandler(useCase *product.UseCase, aiUseCase *product.AIUseCase) *ProductHandler {
	return &ProductHandler{useCase: useCase, aiUseCase: aiUseCase}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	filter := domProduct.ListFilter{
		Limit:  20,
		Offset: 0,
	}

	if catID := r.URL.Query().Get("project_category_id"); catID != "" {
		id, _ := strconv.ParseInt(catID, 10, 64)
		filter.ProjectCategoryID = &id
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := domProduct.ProjectStatus(status)
		filter.Status = &s
	}
	if minP := r.URL.Query().Get("min_price"); minP != "" {
		p, _ := strconv.ParseFloat(minP, 64)
		filter.MinBudget = &p
	}
	if maxP := r.URL.Query().Get("max_price"); maxP != "" {
		p, _ := strconv.ParseFloat(maxP, 64)
		filter.MaxBudget = &p
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		l, _ := strconv.Atoi(limit)
		filter.Limit = l
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		o, _ := strconv.Atoi(offset)
		filter.Offset = o
	}
	filter.Search = r.URL.Query().Get("search")
	filter.SortBy = r.URL.Query().Get("sort_by")
	filter.SortOrder = r.URL.Query().Get("sort_order")

	projects, total, err := h.useCase.ListProjects(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if p == nil {
			respondWithError(w, http.StatusNotFound, "project not found")
			return
		}
		h.useCase.IncrementViews(r.Context(), p.ID)
		respondWithJSON(w, http.StatusOK, dto.FromProject(p))
		return
	}

	p, err := h.useCase.GetProject(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if p == nil {
		respondWithError(w, http.StatusNotFound, "project not found")
		return
	}
	h.useCase.IncrementViews(r.Context(), p.ID)
	respondWithJSON(w, http.StatusOK, dto.FromProject(p))
}

func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.useCase.ListCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	projects, err := h.aiUseCase.Search(r.Context(), req.Query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
		respondWithError(w, http.StatusBadRequest, "Не удалось получить файл")
		return
	}
	defer file.Close()

	if handler.Size > 10<<20 {
		respondWithError(w, http.StatusBadRequest, "Файл слишком большой")
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"url": "/uploads/" + filename,
	})
}
func (h *ProductHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
	projects, _, err := h.useCase.ListProjects(r.Context(), domProduct.ListFilter{Limit: 10000})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
	f.Write(w)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromProject(p))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.CreateProjectRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err := h.useCase.GetProject(r.Context(), id)
	if err != nil || p == nil {
		respondWithError(w, http.StatusNotFound, "project not found")
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dto.FromProject(p))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.useCase.DeleteProject(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
