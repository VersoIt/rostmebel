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

	if catID := r.URL.Query().Get("category_id"); catID != "" {
		id, _ := strconv.ParseInt(catID, 10, 64)
		filter.CategoryID = &id
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := domProduct.ProductStatus(status)
		filter.Status = &s
	}
	if minP := r.URL.Query().Get("min_price"); minP != "" {
		p, _ := strconv.ParseFloat(minP, 64)
		filter.MinPrice = &p
	}
	if maxP := r.URL.Query().Get("max_price"); maxP != "" {
		p, _ := strconv.ParseFloat(maxP, 64)
		filter.MaxPrice = &p
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

	products, total, err := h.useCase.ListProducts(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		res[i] = dto.FromProduct(p)
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
		p, err := h.useCase.GetProductBySlug(r.Context(), idStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if p == nil {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		h.useCase.IncrementViews(r.Context(), p.ID)
		respondWithJSON(w, http.StatusOK, dto.FromProduct(p))
		return
	}

	p, err := h.useCase.GetProduct(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if p == nil {
		respondWithError(w, http.StatusNotFound, "product not found")
		return
	}
	h.useCase.IncrementViews(r.Context(), p.ID)
	respondWithJSON(w, http.StatusOK, dto.FromProduct(p))
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

	products, err := h.aiUseCase.Search(r.Context(), req.Query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		res[i] = dto.FromProduct(p)
	}

	respondWithJSON(w, http.StatusOK, res)
}
func (h *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// r.FormFile automatically calls r.ParseMultipartForm
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Printf("Upload error: %v\n", err) // Log to console for debugging
		respondWithError(w, http.StatusBadRequest, "Не удалось получить файл. Убедитесь, что Content-Type верный.")
		return
	}
	defer file.Close()

	// 10MB limit check
	if handler.Size > 10<<20 {
		respondWithError(w, http.StatusBadRequest, "Файл слишком большой (макс. 10МБ)")
		return
	}

	// Ensure upload dir exists relative to the current working directory
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка создания папки для загрузок")
		return
	}

	// Create unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Не удалось создать файл на сервере")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при сохранении содержимого файла")
		return
	}

	// Absolute URL for the client
	respondWithJSON(w, http.StatusOK, map[string]string{
		"url": "/uploads/" + filename,
	})
}
func (h *ProductHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
	products, _, err := h.useCase.ListProducts(r.Context(), domProduct.ListFilter{Limit: 10000})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Товары"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Название", "Slug", "Цена", "Статус", "Просмотры", "Заказы"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, head)
	}

	for i, p := range products {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Slug)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), string(p.Status))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.ViewsCount)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.OrdersCount)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=products.xlsx")
	f.Write(w)
}

// Admin Handlers
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := &domProduct.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		PriceOld:    req.PriceOld,
		Images:      req.Images,
		Specs:       req.Specs,
		AITags:      req.AITags,
		Status:      req.Status,
	}

	if p.Status == "" {
		p.Status = domProduct.StatusDraft
	}

	if err := h.useCase.CreateProduct(r.Context(), p); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromProduct(p))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.CreateProductRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err := h.useCase.GetProduct(r.Context(), id)
	if err != nil || p == nil {
		respondWithError(w, http.StatusNotFound, "product not found")
		return
	}

	p.CategoryID = req.CategoryID
	p.Name = req.Name
	p.Slug = req.Slug
	p.Description = req.Description
	p.Price = req.Price
	p.PriceOld = req.PriceOld
	p.Images = req.Images
	p.Specs = req.Specs
	p.AITags = req.AITags
	p.Status = req.Status

	if err := h.useCase.UpdateProduct(r.Context(), p); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dto.FromProduct(p))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.useCase.DeleteProduct(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
