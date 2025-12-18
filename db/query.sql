-- name: CreateUser :one
INSERT INTO users (name, dob, email, password_hash, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, name, dob FROM users
ORDER BY id;

-- name: UpdateUser :one
-- UPDATE users
-- SET name = $2, dob = $3
-- WHERE id = $1
-- RETURNING id, name, dob;

UPDATE users
SET name = $2, dob = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

UPDATE users
SET password_hash= $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;