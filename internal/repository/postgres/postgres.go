package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id            SERIAL PRIMARY KEY,
	login         VARCHAR UNIQUE NOT NULL,
	password_hash VARCHAR NOT NULL,
	created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
	id          SERIAL PRIMARY KEY,
	user_id     INTEGER NOT NULL REFERENCES users(id),
	number      VARCHAR UNIQUE NOT NULL,
	status      VARCHAR NOT NULL DEFAULT 'NEW',
	accrual     NUMERIC(15,4),
	uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS balances (
	user_id   INTEGER PRIMARY KEY REFERENCES users(id),
	current   NUMERIC(15,4) NOT NULL DEFAULT 0,
	withdrawn NUMERIC(15,4) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS withdrawals (
	id           SERIAL PRIMARY KEY,
	user_id      INTEGER NOT NULL REFERENCES users(id),
	order_number VARCHAR NOT NULL,
	sum          NUMERIC(15,4) NOT NULL,
	processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if _, err = pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
