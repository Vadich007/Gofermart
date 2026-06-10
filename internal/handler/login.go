package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Vadich007/Gofermart/internal/repository"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil || creds.Login == "" || creds.Password == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := h.users.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
