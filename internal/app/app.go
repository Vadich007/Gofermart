package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

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
		slog.Info("accrual system address not set, worker disabled")
	}

	srv := &http.Server{Addr: cfg.RunAddress, Handler: h.Router()}

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			slog.Error("graceful shutdown failed", "err", err)
		}
	}()

	slog.Info("starting server", "addr", cfg.RunAddress)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
