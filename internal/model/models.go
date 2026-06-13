package model

import "time"

type User struct {
	ID           int
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int
	UserID     int
	Number     string
	Status     OrderStatus
	Accrual    *float64
	UploadedAt time.Time
}

type Balance struct {
	Current   float64
	Withdrawn float64
}

type Withdrawal struct {
	UserID      int
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}
