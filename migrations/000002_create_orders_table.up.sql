CREATE TABLE IF NOT EXISTS orders (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    number      VARCHAR UNIQUE NOT NULL,
    status      VARCHAR NOT NULL DEFAULT 'NEW',
    accrual     NUMERIC(15,4),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
