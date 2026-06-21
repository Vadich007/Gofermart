package handler_test

import (
	"context"
	"iter"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

func sliceToSeq[T any](items []T, err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if err != nil {
			var zero T
			yield(zero, err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

type mockUserRepo struct {
	createFn     func(ctx context.Context, login, passwordHash string) (*model.User, error)
	getByLoginFn func(ctx context.Context, login string) (*model.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, login, passwordHash string) (*model.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, login, passwordHash)
	}
	return &model.User{ID: 1, Login: login, PasswordHash: passwordHash}, nil
}

func (m *mockUserRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	if m.getByLoginFn != nil {
		return m.getByLoginFn(ctx, login)
	}
	return nil, repository.ErrNotFound
}

type mockOrderRepo struct {
	createFn       func(ctx context.Context, userID int, number string) error
	getByUserFn    func(ctx context.Context, userID int) ([]*model.Order, error)
	getPendingFn   func(ctx context.Context) ([]*model.Order, error)
	updateStatusFn func(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
}

func (m *mockOrderRepo) Create(ctx context.Context, userID int, number string) error {
	if m.createFn != nil {
		return m.createFn(ctx, userID, number)
	}
	return nil
}

func (m *mockOrderRepo) GetByUser(ctx context.Context, userID int) iter.Seq2[*model.Order, error] {
	if m.getByUserFn != nil {
		orders, err := m.getByUserFn(ctx, userID)
		return sliceToSeq(orders, err)
	}
	return sliceToSeq[*model.Order](nil, nil)
}

func (m *mockOrderRepo) GetPending(ctx context.Context) iter.Seq2[*model.Order, error] {
	if m.getPendingFn != nil {
		orders, err := m.getPendingFn(ctx)
		return sliceToSeq(orders, err)
	}
	return sliceToSeq[*model.Order](nil, nil)
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, number, status, accrual)
	}
	return nil
}

type mockBalanceRepo struct {
	initFn           func(ctx context.Context, userID int) error
	getFn            func(ctx context.Context, userID int) (*model.Balance, error)
	addAccrualFn     func(ctx context.Context, userID int, amount float64) error
	withdrawFn       func(ctx context.Context, userID int, orderNumber string, sum float64) error
	getWithdrawalsFn func(ctx context.Context, userID int) ([]*model.Withdrawal, error)
}

func (m *mockBalanceRepo) Init(ctx context.Context, userID int) error {
	if m.initFn != nil {
		return m.initFn(ctx, userID)
	}
	return nil
}

func (m *mockBalanceRepo) Get(ctx context.Context, userID int) (*model.Balance, error) {
	if m.getFn != nil {
		return m.getFn(ctx, userID)
	}
	return &model.Balance{}, nil
}

func (m *mockBalanceRepo) AddAccrual(ctx context.Context, userID int, amount float64) error {
	if m.addAccrualFn != nil {
		return m.addAccrualFn(ctx, userID, amount)
	}
	return nil
}

func (m *mockBalanceRepo) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	if m.withdrawFn != nil {
		return m.withdrawFn(ctx, userID, orderNumber, sum)
	}
	return nil
}

func (m *mockBalanceRepo) GetWithdrawals(ctx context.Context, userID int) iter.Seq2[*model.Withdrawal, error] {
	if m.getWithdrawalsFn != nil {
		withdrawals, err := m.getWithdrawalsFn(ctx, userID)
		return sliceToSeq(withdrawals, err)
	}
	return sliceToSeq[*model.Withdrawal](nil, nil)
}
