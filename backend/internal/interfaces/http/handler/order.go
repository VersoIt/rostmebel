package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rostmebel/backend/internal/application/order"
	"github.com/rostmebel/backend/internal/domain/apperror"
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
		respondWithError(w, err)
		return
	}

	// Honeypot check
	if req.Website != "" {
		w.WriteHeader(http.StatusCreated)
		return
	}

	o := &domOrder.Order{
		ProjectID:     req.ProjectID,
		ClientName:    req.ClientName,
		ClientPhone:   req.ClientPhone,
		ClientEmail:   req.ClientEmail,
		Comment:       req.Comment,
		ProjectType:   req.ProjectType,
		BudgetRange:   req.BudgetRange,
		City:          req.City,
		ContactMethod: req.ContactMethod,
		Status:        domOrder.StatusNew,
		IPAddress:     clientIP(r),
		UserAgent:     r.UserAgent(),
		Fingerprint:   req.Fingerprint,
	}

	if err := h.useCase.CreateOrder(r.Context(), o); err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FromOrder(o))
}

func clientIP(r *http.Request) net.IP {
	candidates := make([]string, 0, 3)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		candidates = append(candidates, strings.Split(forwardedFor, ",")...)
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		candidates = append(candidates, realIP)
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		candidates = append(candidates, host)
	} else {
		candidates = append(candidates, r.RemoteAddr)
	}

	for _, candidate := range candidates {
		ip := net.ParseIP(strings.TrimSpace(candidate))
		if ip != nil {
			return ip
		}
	}
	return net.IPv4(0, 0, 0, 0)
}

// Admin Handlers
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	filter := domOrder.ListFilter{
		Status: domOrder.OrderStatus(r.URL.Query().Get("status")),
	}
	if filter.Status != "" && !isKnownOrderStatus(filter.Status) {
		respondWithError(w, invalidQuery("status", string(filter.Status)))
		return
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 1 || l > 100 {
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

	orders, filteredTotal, absoluteTotal, err := h.useCase.ListOrders(r.Context(), filter)
	if err != nil {
		respondWithError(w, err)
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, apperror.New(apperror.CodeInvalidID, "Invalid order id", map[string]any{"id": idStr}))
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := decodeAndValidate(r, &req); err != nil {
		respondWithError(w, err)
		return
	}
	if !isKnownOrderStatus(req.Status) {
		respondWithError(w, invalidQuery("status", string(req.Status)))
		return
	}

	if err := h.useCase.UpdateOrderStatus(r.Context(), id, req.Status); err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) MarkAsSpam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, apperror.New(apperror.CodeInvalidID, "Invalid order id", map[string]any{"id": idStr}))
		return
	}

	if err := h.useCase.MarkAsSpam(r.Context(), id); err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.useCase.ExportOrders(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Заявки"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"ID", "Имя клиента", "Телефон", "Email", "Тип проекта", "Бюджет", "Город", "Связь", "Комментарий", "Статус", "Дата"}
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
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), o.ProjectType)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), o.BudgetRange)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), o.City)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), o.ContactMethod)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), o.Comment)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), string(o.Status))
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), o.CreatedAt.Format("2006-01-02 15:04"))
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=orders.xlsx")

	if err := f.Write(w); err != nil {
		respondWithError(w, err)
	}
}

func isKnownOrderStatus(status domOrder.OrderStatus) bool {
	switch status {
	case domOrder.StatusNew, domOrder.StatusProcessing, domOrder.StatusDone, domOrder.StatusRejected, domOrder.StatusSpam:
		return true
	default:
		return false
	}
}
