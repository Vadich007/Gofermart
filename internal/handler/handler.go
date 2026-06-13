package handler

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/Vadich007/Gofermart/internal/middleware"
	"github.com/Vadich007/Gofermart/internal/service"
)

type Handler struct {
	users    *service.UserService
	orders   *service.OrderService
	balances *service.BalanceService
}

func New(users *service.UserService, orders *service.OrderService, balances *service.BalanceService) *Handler {
	return &Handler{users: users, orders: orders, balances: balances}
}

func (h *Handler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.Compress)

	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(mw.Auth)
		r.Post("/api/user/orders", h.UploadOrder)
		r.Get("/api/user/orders", h.GetOrders)
		r.Get("/api/user/balance", h.GetBalance)
		r.Post("/api/user/balance/withdraw", h.Withdraw)
		r.Get("/api/user/withdrawals", h.GetWithdrawals)
	})

	return r
}
