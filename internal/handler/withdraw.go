package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	mw "github.com/Vadich007/Gofermart/internal/middleware"
	"github.com/Vadich007/Gofermart/internal/repository"
	"github.com/Vadich007/Gofermart/internal/service"
)

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req withdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Order == "" || req.Sum <= 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID := mw.GetUserID(r.Context())

	err := h.balances.Withdraw(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrderNumber):
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		case errors.Is(err, repository.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
