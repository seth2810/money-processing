-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    id serial PRIMARY KEY,
    email varchar(320) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE clients CASCADE;
DROP TABLE clients;
-- +goose StatementEnd
