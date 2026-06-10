package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	router := buildHandler(
		&mockUserRepo{
			getByLoginFn: func(_ context.Context, login string) (*model.User, error) {
				return &model.User{ID: 1, Login: login, PasswordHash: string(hash)}, nil
			},
		},
		&mockOrderRepo{},
		&mockBalanceRepo{},
	)

	body := `{"login":"alice","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := do(router, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("Authorization") == "" {
		t.Error("expected Authorization header in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	router := buildHandler(
		&mockUserRepo{
			getByLoginFn: func(_ context.Context, _ string) (*model.User, error) {
				return &model.User{ID: 1, PasswordHash: string(hash)}, nil
			},
		},
		&mockOrderRepo{},
		&mockBalanceRepo{},
	)

	body := `{"login":"alice","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := do(router, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	router := buildHandler(
		&mockUserRepo{
			getByLoginFn: func(_ context.Context, _ string) (*model.User, error) {
				return nil, repository.ErrNotFound
			},
		},
		&mockOrderRepo{},
		&mockBalanceRepo{},
	)

	body := `{"login":"nobody","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := do(router, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestLogin_BadJSON(t *testing.T) {
	router := buildHandler(&mockUserRepo{}, &mockOrderRepo{}, &mockBalanceRepo{})

	cases := []string{`not json`, `{}`, `{"login":""}`}
	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := do(router, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("body=%q: status = %d, want %d", body, rr.Code, http.StatusBadRequest)
		}
	}
}
