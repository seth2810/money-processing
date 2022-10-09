-- +goose Up
-- +goose StatementBegin
CREATE TYPE currency_ticker AS ENUM ('USD', 'COP', 'MXN');

CREATE TABLE accounts (
    id serial PRIMARY KEY,
    client_id integer NOT NULL REFERENCES clients(id),
    balance decimal(10, 4) NOT NULL CHECK (balance >= 0) DEFAULT 0,
    currency currency_ticker NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE accounts CASCADE;
DROP TABLE accounts;
DROP TYPE currency_ticker;
-- +goose StatementEnd
