package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

func TestRegister_Success(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{
			createFn: func(_ context.Context, login, _ string) (*model.User, error) {
				return &model.User{ID: 1, Login: login}, nil
			},
		},
		&mockOrderRepo{},
		&mockBalanceRepo{},
	)

	body := `{"login":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("Authorization") == "" {
		t.Error("expected Authorization header in response")
	}
}

func TestRegister_BadJSON(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	cases := []string{
		"not json",
		`{"login":""}`,
		`{"password":"pass"}`,
		`{}`,
	}

	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := do(router, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("body=%q: status = %d, want %d", body, rr.Code, http.StatusBadRequest)
		}
	}
}

func TestRegister_LoginConflict(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{
			createFn: func(_ context.Context, _, _ string) (*model.User, error) {
				return nil, repository.ErrLoginConflict
			},
		},
		&mockOrderRepo{},
		&mockBalanceRepo{},
	)

	body := `{"login":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := do(router, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusConflict)
	}
}
