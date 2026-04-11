package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/product"
	domProduct "github.com/rostmebel/backend/internal/domain/product"
	"github.com/rostmebel/backend/internal/interfaces/dto"
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
