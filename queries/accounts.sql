-- name: CreateAccount :one
INSERT INTO accounts (client_id, currency) VALUES (@client_id, @currency) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = @id LIMIT 1;

-- name: DepositToAccount :exec
UPDATE accounts SET balance = balance + @amount WHERE id = @account_id;

-- name: WithdrawFromAccount :exec
UPDATE accounts SET balance = balance - @amount WHERE id = @account_id;
