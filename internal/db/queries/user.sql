-- name: GetAllUsers :many
SELECT id, name, email, profile_img, created_at, updated_at FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: CountUserByEmail :one
SELECT COUNT(1) FROM users WHERE email = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: CreateUser :execresult
INSERT INTO users (id, name, email, password)
VALUES (?, ?, ?, ?);