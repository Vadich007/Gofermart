CREATE TABLE IF NOT EXISTS withdrawals (
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES users(id),
    order_number VARCHAR NOT NULL,
    sum          NUMERIC(15,4) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
