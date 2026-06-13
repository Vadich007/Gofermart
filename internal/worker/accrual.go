package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

type AccrualWorker struct {
	orders   repository.OrderRepository
	balances repository.BalanceRepository
	baseURL  string
	client   *http.Client
	interval time.Duration
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func New(
	orders repository.OrderRepository,
	balances repository.BalanceRepository,
	accrualBaseURL string,
) *AccrualWorker {
	return &AccrualWorker{
		orders:   orders,
		balances: balances,
		baseURL:  accrualBaseURL,
		client:   &http.Client{Timeout: 10 * time.Second},
		interval: 2 * time.Second,
	}
}

func (w *AccrualWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processOrders(ctx)
		}
	}
}

func (w *AccrualWorker) processOrders(ctx context.Context) {
	for order, err := range w.orders.GetPending(ctx) {
		if err != nil {
			slog.Error("accrual worker: get pending orders", "err", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		retryAfter, err := w.processOrder(ctx, order)
		if err != nil {
			slog.Error("accrual worker: process order", "order", order.Number, "err", err)
		}
		if retryAfter > 0 {
			slog.Info("accrual worker: rate limited, sleeping", "duration", retryAfter)
			select {
			case <-ctx.Done():
				return
			case <-time.After(retryAfter):
			}
			return
		}
	}
}

func (w *AccrualWorker) processOrder(ctx context.Context, order *model.Order) (time.Duration, error) {
	url := fmt.Sprintf("%s/api/orders/%s", w.baseURL, order.Number)
	resp, err := w.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		retryAfter := 60 * time.Second
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return retryAfter, nil

	case http.StatusNoContent:
		return 0, nil

	case http.StatusInternalServerError:
		return 0, fmt.Errorf("accrual service error 500 for order %s", order.Number)

	case http.StatusOK:
		var ar accrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
			return 0, err
		}
		return 0, w.applyAccrual(ctx, order, &ar)
	}

	return 0, nil
}

func (w *AccrualWorker) applyAccrual(ctx context.Context, order *model.Order, ar *accrualResponse) error {
	var newStatus model.OrderStatus
	var accrual *float64

	switch ar.Status {
	case "REGISTERED":
		newStatus = model.OrderStatusNew
	case "PROCESSING":
		newStatus = model.OrderStatusProcessing
	case "INVALID":
		newStatus = model.OrderStatusInvalid
	case "PROCESSED":
		newStatus = model.OrderStatusProcessed
		accrual = &ar.Accrual
	default:
		return nil
	}

	if err := w.orders.UpdateStatus(ctx, order.Number, newStatus, accrual); err != nil {
		return err
	}

	if newStatus == model.OrderStatusProcessed && accrual != nil && *accrual > 0 {
		if err := w.balances.AddAccrual(ctx, order.UserID, *accrual); err != nil {
			return err
		}
	}

	return nil
}
