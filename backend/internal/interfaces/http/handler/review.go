package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/review"
	domReview "github.com/rostmebel/backend/internal/domain/review"
	"github.com/rostmebel/backend/internal/interfaces/dto"
)

type ReviewHandler struct {
	useCase *review.UseCase
}

func NewReviewHandler(useCase *review.UseCase) *ReviewHandler {
	return &ReviewHandler{useCase: useCase}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReviewRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	rev := &domReview.Review{
		ProjectID: req.ProjectID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Images:    req.Images,
	}

	if err := h.useCase.CreateReview(r.Context(), req.ClientPhone, rev); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromReview(rev))
}

func (h *ReviewHandler) GetProjectReviews(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "id")
	projectID, _ := strconv.ParseInt(projectIDStr, 10, 64)

	reviews, err := h.useCase.GetByProject(r.Context(), projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.ReviewResponse, len(reviews))
	for i, rev := range reviews {
		res[i] = dto.FromReview(rev)
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (h *ReviewHandler) AdminListReviews(w http.ResponseWriter, r *http.Request) {
	filter := domReview.ListFilter{
		Status: domReview.ReviewStatus(r.URL.Query().Get("status")),
		Limit:  20,
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		filter.Limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		filter.Offset, _ = strconv.Atoi(o)
	}

	reviews, filteredTotal, absoluteTotal, err := h.useCase.ListReviews(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.ReviewResponse, len(reviews))
	for i, rev := range reviews {
		res[i] = dto.FromReview(rev)
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items":          res,
		"total":          filteredTotal,
		"absolute_total": absoluteTotal,
	})
}

func (h *ReviewHandler) AdminModerateReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.ModerateReviewRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.ModerateReview(r.Context(), id, req.Approved); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReviewHandler) AdminDeleteReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.useCase.DeleteReview(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
