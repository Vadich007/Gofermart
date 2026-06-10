package service

import (
	"context"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

type BalanceService struct {
	balances repository.BalanceRepository
}

func NewBalanceService(balances repository.BalanceRepository) *BalanceService {
	return &BalanceService{balances: balances}
}

func (s *BalanceService) Get(ctx context.Context, userID int) (*model.Balance, error) {
	return s.balances.Get(ctx, userID)
}

func (s *BalanceService) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	if !isValidLuhn(orderNumber) {
		return ErrInvalidOrderNumber
	}
	return s.balances.Withdraw(ctx, userID, orderNumber, sum)
}

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdrawal, error) {
	return s.balances.GetWithdrawals(ctx, userID)
}
