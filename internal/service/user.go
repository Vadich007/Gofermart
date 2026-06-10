package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Vadich007/Gofermart/internal/repository"
)

const jwtSecret = "gophermart-secret-key"
const tokenTTL = 24 * time.Hour

type UserService struct {
	users    repository.UserRepository
	balances repository.BalanceRepository
}

func NewUserService(users repository.UserRepository, balances repository.BalanceRepository) *UserService {
	return &UserService{users: users, balances: balances}
}

func (s *UserService) Register(ctx context.Context, login, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := s.users.Create(ctx, login, string(hash))
	if err != nil {
		return "", err
	}

	if err := s.balances.Init(ctx, user.ID); err != nil {
		return "", err
	}

	return buildToken(user.ID)
}

func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", repository.ErrNotFound
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", repository.ErrNotFound
	}

	return buildToken(user.ID)
}

func buildToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id claim")
	}

	return int(uid), nil
}
