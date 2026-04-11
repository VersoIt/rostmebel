package handler

import (
	"net/http"

	"github.com/rostmebel/backend/internal/application/admin"
	"github.com/rostmebel/backend/internal/interfaces/dto"
)

type AdminHandler struct {
	useCase *admin.UseCase
}

func NewAdminHandler(useCase *admin.UseCase) *AdminHandler {
	return &AdminHandler{useCase: useCase}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.useCase.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *AdminHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.useCase.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// sub is set in JWT middleware
	adminID := r.Context().Value("sub").(int64)
	if err := h.useCase.Logout(r.Context(), adminID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.useCase.GetStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}
