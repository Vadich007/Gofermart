package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, userID int, number string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO orders (user_id, number) VALUES ($1, $2)`,
		userID, number,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			var ownerID int
			_ = r.db.QueryRow(ctx,
				`SELECT user_id FROM orders WHERE number = $1`, number,
			).Scan(&ownerID)
			if ownerID == userID {
				return repository.ErrOrderOwned
			}
			return repository.ErrOrderConflict
		}
		return err
	}
	return nil
}

func (r *OrderRepo) GetByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, number, status, accrual, uploaded_at
		 FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		o := &model.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *OrderRepo) GetPending(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, number, status, accrual, uploaded_at
		 FROM orders WHERE status IN ('NEW', 'PROCESSING')`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		o := &model.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`,
		status, accrual, number,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
