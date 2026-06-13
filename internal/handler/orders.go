package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	mw "github.com/Vadich007/Gofermart/internal/middleware"
	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
	"github.com/Vadich007/Gofermart/internal/service"
)

type orderResponse struct {
	Number     string      `json:"number"`
	Status     string      `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (h *Handler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	number := string(body)
	userID := mw.GetUserID(r.Context())

	err = h.orders.Upload(r.Context(), userID, number)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrderNumber):
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		case errors.Is(err, repository.ErrOrderOwned):
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, repository.ErrOrderConflict):
			http.Error(w, "order uploaded by another user", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r.Context())

	orders, err := h.orders.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]orderResponse, len(orders))
	for i, o := range orders {
		resp[i] = toOrderResponse(o)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

func toOrderResponse(o *model.Order) orderResponse {
	return orderResponse{
		Number:     o.Number,
		Status:     string(o.Status),
		Accrual:    o.Accrual,
		UploadedAt: o.UploadedAt,
	}
}
