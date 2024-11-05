-- name: CreateUser :one
INSERT INTO users (name, email, address, user_type, password_hash, profile_headline)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING name, email, user_type;

-- name: GetUser :one
SELECT id, password_hash, user_type FROM users
WHERE email = $1;
