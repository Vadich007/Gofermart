package service

import (
	"context"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

type OrderService struct {
	orders repository.OrderRepository
}

func NewOrderService(orders repository.OrderRepository) *OrderService {
	return &OrderService{orders: orders}
}

func (s *OrderService) Upload(ctx context.Context, userID int, number string) error {
	if !isValidLuhn(number) {
		return ErrInvalidOrderNumber
	}
	return s.orders.Create(ctx, userID, number)
}

func (s *OrderService) GetByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	var orders []*model.Order
	for o, err := range s.orders.GetByUser(ctx, userID) {
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func isValidLuhn(number string) bool {
	if len(number) == 0 {
		return false
	}
	sum := 0
	nDigits := len(number)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		ch := number[i]
		if ch < '0' || ch > '9' {
			return false
		}
		digit := int(ch - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
