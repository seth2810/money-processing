-- name: ListTransactions :many
SELECT *
FROM transactions
WHERE from_account_id = @account_id::integer
    OR to_account_id = @account_id::integer;

-- name: CreateDepositTransaction :one
INSERT INTO transactions (type, amount, to_account_id) VALUES ('deposit', @amount, @account_id::integer) RETURNING *;

-- name: CreateWithdrawTransaction :one
INSERT INTO transactions (type, amount, from_account_id) VALUES ('withdraw', @amount, @account_id::integer) RETURNING *;

-- name: CreateTransferTransaction :one
INSERT INTO transactions (type, amount, from_account_id, to_account_id) VALUES ('transfer', @amount, @from_account_id::integer, @to_account_id::integer) RETURNING *;
