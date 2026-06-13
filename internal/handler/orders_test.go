package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

func TestUploadOrder_Accepted(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903"))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
}

func TestUploadOrder_AlreadyOwnedByThisUser(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{
			createFn: func(_ context.Context, _ int, _ string) error {
				return repository.ErrOrderOwned
			},
		},
		&mockBalanceRepo{},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903"))
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestUploadOrder_ConflictOtherUser(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{
			createFn: func(_ context.Context, _ int, _ string) error {
				return repository.ErrOrderConflict
			},
		},
		&mockBalanceRepo{},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903"))
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusConflict)
	}
}

func TestUploadOrder_InvalidLuhn(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678900"))
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnprocessableEntity)
	}
}

func TestUploadOrder_Unauthorized(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903"))

	rr := do(router, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestGetOrders_WithData(t *testing.T) {
	accrual := 500.0
	orders := []*model.Order{
		{
			Number:     "12345678903",
			Status:     model.OrderStatusProcessed,
			Accrual:    &accrual,
			UploadedAt: time.Now(),
		},
	}

	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{
			getByUserFn: func(_ context.Context, _ int) ([]*model.Order, error) {
				return orders, nil
			},
		},
		&mockBalanceRepo{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("got %d orders, want 1", len(result))
	}
	if result[0]["number"] != "12345678903" {
		t.Errorf("order number = %v, want 12345678903", result[0]["number"])
	}
}

func TestGetOrders_NoData(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{
			getByUserFn: func(_ context.Context, _ int) ([]*model.Order, error) {
				return nil, nil
			},
		},
		&mockBalanceRepo{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}
