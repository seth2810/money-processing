-- name: CreateClient :one
INSERT INTO clients (email) VALUES (@email) RETURNING *;

-- name: GetClient :one
SELECT * FROM clients WHERE id = @id LIMIT 1;
