package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/order"
	domOrder "github.com/rostmebel/backend/internal/domain/order"
	"github.com/rostmebel/backend/internal/interfaces/dto"
	"github.com/xuri/excelize/v2"
)

type OrderHandler struct {
	useCase *order.UseCase
}

func NewOrderHandler(useCase *order.UseCase) *OrderHandler {
	return &OrderHandler{useCase: useCase}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Honeypot check
	if req.Website != "" {
		w.WriteHeader(http.StatusCreated)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if forward := r.Header.Get("X-Forwarded-For"); forward != "" {
		ip = forward
	}

	o := &domOrder.Order{
		ProjectID:   req.ProjectID,
		ClientName:  req.ClientName,
		ClientPhone: req.ClientPhone,
		ClientEmail: req.ClientEmail,
		Comment:     req.Comment,
		Status:      domOrder.StatusNew,
		IPAddress:   net.ParseIP(ip),
		UserAgent:   r.UserAgent(),
		Fingerprint: req.Fingerprint,
	}

	if err := h.useCase.CreateOrder(r.Context(), o); err != nil {
		respondWithError(w, http.StatusTooManyRequests, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromOrder(o))
}

// Admin Handlers
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	filter := domOrder.ListFilter{
		Status: domOrder.OrderStatus(r.URL.Query().Get("status")),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}

	orders, filteredTotal, absoluteTotal, err := h.useCase.ListOrders(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		res[i] = dto.FromOrder(o)
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items":          res,
		"total":          filteredTotal,
		"absolute_total": absoluteTotal,
	})
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.UpdateOrderStatusRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.UpdateOrderStatus(r.Context(), id, req.Status); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) MarkAsSpam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.useCase.MarkAsSpam(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.useCase.ExportOrders(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Заявки"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"ID", "Имя клиента", "Телефон", "Email", "Комментарий", "Статус", "Дата"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, head)
	}

	// Data
	for i, o := range orders {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), o.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), o.ClientName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), o.ClientPhone)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), o.ClientEmail)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), o.Comment)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), string(o.Status))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), o.CreatedAt.Format("2006-01-02 15:04"))
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=orders.xlsx")

	if err := f.Write(w); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
