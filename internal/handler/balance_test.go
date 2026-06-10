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

func TestGetBalance_Success(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{},
		&mockBalanceRepo{
			getFn: func(_ context.Context, _ int) (*model.Balance, error) {
				return &model.Balance{Current: 500.5, Withdrawn: 42}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var result map[string]float64
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result["current"] != 500.5 {
		t.Errorf("current = %v, want 500.5", result["current"])
	}
	if result["withdrawn"] != 42 {
		t.Errorf("withdrawn = %v, want 42", result["withdrawn"])
	}
}

func TestGetBalance_Unauthorized(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	rr := do(router, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestWithdraw_Success(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{},
		&mockBalanceRepo{
			withdrawFn: func(_ context.Context, _ int, _ string, _ float64) error {
				return nil
			},
		},
	)

	body := `{"order":"2377225624","sum":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{},
		&mockBalanceRepo{
			withdrawFn: func(_ context.Context, _ int, _ string, _ float64) error {
				return repository.ErrInsufficientFunds
			},
		},
	)

	body := `{"order":"2377225624","sum":9999}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusPaymentRequired {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusPaymentRequired)
	}
}

func TestWithdraw_InvalidOrder(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	body := `{"order":"12345678900","sum":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnprocessableEntity)
	}
}

func TestGetWithdrawals_WithData(t *testing.T) {
	withdrawals := []*model.Withdrawal{
		{OrderNumber: "2377225624", Sum: 500, ProcessedAt: time.Now()},
	}

	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{},
		&mockBalanceRepo{
			getWithdrawalsFn: func(_ context.Context, _ int) ([]*model.Withdrawal, error) {
				return withdrawals, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
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
		t.Errorf("got %d withdrawals, want 1", len(result))
	}
	if result[0]["order"] != "2377225624" {
		t.Errorf("order = %v, want 2377225624", result[0]["order"])
	}
}

func TestGetWithdrawals_NoData(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{},
		&mockOrderRepo{},
		&mockBalanceRepo{
			getWithdrawalsFn: func(_ context.Context, _ int) ([]*model.Withdrawal, error) {
				return nil, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req.Header.Set("Authorization", "Bearer "+makeToken(t, 1))

	rr := do(router, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}
