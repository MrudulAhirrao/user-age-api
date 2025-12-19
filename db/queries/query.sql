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
SET 
  -- We use sqlc.arg(name) to force the Go struct field to be "Name"
  name = COALESCE(NULLIF(sqlc.arg(name), ''), name), 
  
  -- We use sqlc.arg(dob) to force the Go struct field to be "Dob"
  dob = COALESCE(sqlc.arg(dob), dob),
  
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

UPDATE users
SET password_hash= $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;