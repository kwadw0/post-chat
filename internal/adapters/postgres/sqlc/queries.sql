-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: GetProduct :one
SELECT * FROM products
WHERE uuid = $1 LIMIT 1;

-- name: CreateProduct :one
INSERT INTO products (uuid, name, description, price, quantity) VALUES (gen_random_uuid(), $1, $2, $3, $4) RETURNING *;