package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Vadich007/Gofermart/internal/config"
	"github.com/Vadich007/Gofermart/internal/handler"
	"github.com/Vadich007/Gofermart/internal/repository/postgres"
	"github.com/Vadich007/Gofermart/internal/service"
	"github.com/Vadich007/Gofermart/internal/worker"
)

func Run(ctx context.Context, cfg *config.Config) error {
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURI)
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepo(pool)
	orderRepo := postgres.NewOrderRepo(pool)
	balanceRepo := postgres.NewBalanceRepo(pool)

	userSvc := service.NewUserService(userRepo, balanceRepo)
	orderSvc := service.NewOrderService(orderRepo)
	balanceSvc := service.NewBalanceService(balanceRepo)

	h := handler.New(userSvc, orderSvc, balanceSvc)

	if cfg.AccrualSystemAddress != "" {
		w := worker.New(orderRepo, balanceRepo, cfg.AccrualSystemAddress)
		go w.Run(ctx)
	} else {
		log.Println("accrual system address not set, worker disabled")
	}

	log.Printf("starting server on %s", cfg.RunAddress)
	return http.ListenAndServe(cfg.RunAddress, h.Router())
}
