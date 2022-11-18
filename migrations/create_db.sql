CREATE TABLE IF NOT EXISTS balance(
    user_id SERIAL PRIMARY KEY,
    balance INTEGER NOT NULL,
    hold INTEGER DEFAULT 0
);

CREATE TYPE transaction_type as ENUM ('transfer', 'capture', 'cancel', 'write-off');

CREATE TABLE IF NOT EXISTS transaction(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type transaction_type NOT NULL,
    service_id INTEGER,
    order_id INTEGER,
    total INTEGER NOT NULL
);