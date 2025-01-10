-- name: GetAllUsers :many
SELECT id, name, email, profile_img, created_at, updated_at FROM users;