package handler

import (
	"encoding/json"
	"net/http"

	mw "github.com/Vadich007/Gofermart/internal/middleware"
)

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r.Context())

	balance, err := h.balances.Get(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balanceResponse{ //nolint:errcheck
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}
