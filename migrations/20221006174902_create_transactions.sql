-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdraw', 'transfer');

CREATE TABLE transactions (
    id serial PRIMARY KEY,
    type transaction_type NOT NULL,
    amount decimal(10, 4) NOT NULL CHECK (amount >= 0),
    from_account_id integer REFERENCES accounts(id) CHECK (from_account_id != to_account_id),
    to_account_id integer REFERENCES accounts(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX transactions_from_account_id_to_account_id_idx ON transactions(from_account_id,to_account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE transactions CASCADE;
DROP TABLE transactions;
DROP TYPE transaction_type;
-- +goose StatementEnd
