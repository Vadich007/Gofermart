package handler

import (
	"encoding/json"
	"net/http"
	"time"

	mw "github.com/Vadich007/Gofermart/internal/middleware"
	"github.com/Vadich007/Gofermart/internal/model"
)

type withdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r.Context())

	withdrawals, err := h.balances.GetWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]withdrawalResponse, len(withdrawals))
	for i, wr := range withdrawals {
		resp[i] = toWithdrawalResponse(wr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

func toWithdrawalResponse(w *model.Withdrawal) withdrawalResponse {
	return withdrawalResponse{
		Order:       w.OrderNumber,
		Sum:         w.Sum,
		ProcessedAt: w.ProcessedAt,
	}
}
