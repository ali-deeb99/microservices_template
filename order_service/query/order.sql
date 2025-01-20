-- name: CreateOrder :one
INSERT INTO orders (name, note, status) 
VALUES ($1, $2, $3)RETURNING name;

-- name: UpdateOrder :exec
UPDATE orders 
SET status = $1 
WHERE id = $2;
