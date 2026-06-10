package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Vadich007/Gofermart/internal/handler"
	"github.com/Vadich007/Gofermart/internal/service"
)

const testJWTSecret = "gophermart-secret-key"

func makeToken(t *testing.T, userID int) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("makeToken: %v", err)
	}
	return s
}

func buildHandler(userRepo *mockUserRepo, orderRepo *mockOrderRepo, balanceRepo *mockBalanceRepo) http.Handler {
	userSvc := service.NewUserService(userRepo, balanceRepo)
	orderSvc := service.NewOrderService(orderRepo)
	balanceSvc := service.NewBalanceService(balanceRepo)
	h := handler.New(userSvc, orderSvc, balanceSvc)
	return h.Router()
}

func do(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
