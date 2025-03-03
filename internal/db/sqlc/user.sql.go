// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package sqlc

import (
	"context"
	"database/sql"
)

const countUserByEmail = `-- name: CountUserByEmail :one
SELECT COUNT(1) FROM users WHERE email = ?
`

func (q *Queries) CountUserByEmail(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUserByEmail, email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (id, name, email, password)
VALUES (?, ?, ?, ?)
`

type CreateUserParams struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, name, email, profile_img, created_at, updated_at FROM users
`

type GetAllUsersRow struct {
	ID         string
	Name       string
	Email      string
	ProfileImg sql.NullString
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.ProfileImg,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, profile_img, created_at, updated_at FROM users WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfileImg,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, password, profile_img, created_at, updated_at FROM users
WHERE id = ?
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.ProfileImg,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
