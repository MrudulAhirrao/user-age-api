-- name: CreateOrUpdateAddress :one
INSERT INTO addresses (
  user_id, line1, line2, city, state, postal_code, country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (user_id) 
DO UPDATE SET 
  line1 = EXCLUDED.line1,
  line2 = EXCLUDED.line2,
  city = EXCLUDED.city,
  state = EXCLUDED.state,
  postal_code = EXCLUDED.postal_code,
  country = EXCLUDED.country,
  updated_at = NOW()
RETURNING *;

-- name: GetAddressByUserID :one
SELECT * FROM addresses WHERE user_id = $1;