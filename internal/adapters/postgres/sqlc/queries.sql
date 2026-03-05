-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: GetProduct :one
SELECT * FROM products
WHERE uuid = $1 LIMIT 1;