package postgres

import (
	"context"
	"errors"
	"iter"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Vadich007/Gofermart/internal/model"
	"github.com/Vadich007/Gofermart/internal/repository"
)

type BalanceRepo struct {
	db *pgxpool.Pool
}

func NewBalanceRepo(db *pgxpool.Pool) *BalanceRepo {
	return &BalanceRepo{db: db}
}

func (r *BalanceRepo) Init(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO balances (user_id, current, withdrawn) VALUES ($1, 0, 0) ON CONFLICT DO NOTHING`,
		userID,
	)
	return err
}

func (r *BalanceRepo) Get(ctx context.Context, userID int) (*model.Balance, error) {
	b := &model.Balance{}
	err := r.db.QueryRow(ctx,
		`SELECT current, withdrawn FROM balances WHERE user_id = $1`,
		userID,
	).Scan(&b.Current, &b.Withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Balance{}, nil
		}
		return nil, err
	}
	return b, nil
}

func (r *BalanceRepo) AddAccrual(ctx context.Context, userID int, amount float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE balances SET current = current + $1 WHERE user_id = $2`,
		amount, userID,
	)
	return err
}

func (r *BalanceRepo) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var current float64
	if err := tx.QueryRow(ctx,
		`SELECT current FROM balances WHERE user_id = $1 FOR UPDATE`,
		userID,
	).Scan(&current); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrInsufficientFunds
		}
		return err
	}

	if current < sum {
		return repository.ErrInsufficientFunds
	}

	if _, err := tx.Exec(ctx,
		`UPDATE balances SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`,
		sum, userID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`,
		userID, orderNumber, sum,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *BalanceRepo) GetWithdrawals(ctx context.Context, userID int) iter.Seq2[*model.Withdrawal, error] {
	return func(yield func(*model.Withdrawal, error) bool) {
		rows, err := r.db.Query(ctx,
			`SELECT user_id, order_number, sum, processed_at
			 FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC`,
			userID,
		)
		if err != nil {
			yield(nil, err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			w := &model.Withdrawal{}
			if err := rows.Scan(&w.UserID, &w.OrderNumber, &w.Sum, &w.ProcessedAt); err != nil {
				yield(nil, err)
				return
			}
			if !yield(w, nil) {
				return
			}
		}
		if err := rows.Err(); err != nil {
			yield(nil, err)
		}
	}
}
