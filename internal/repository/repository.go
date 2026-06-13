package repository

import (
	"context"
	"errors"

	"github.com/Vadich007/Gofermart/internal/model"
)

var (
	ErrLoginConflict     = errors.New("login already taken")
	ErrOrderConflict     = errors.New("order already uploaded by another user")
	ErrOrderOwned        = errors.New("order already uploaded by this user")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrNotFound          = errors.New("not found")
)

type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, userID int, number string) error
	GetByUser(ctx context.Context, userID int) ([]*model.Order, error)
	GetPending(ctx context.Context) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
}

type BalanceRepository interface {
	Init(ctx context.Context, userID int) error
	Get(ctx context.Context, userID int) (*model.Balance, error)
	AddAccrual(ctx context.Context, userID int, amount float64) error
	Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdrawal, error)
}
