// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package sqlc

import (
	"database/sql"
)

type User struct {
	ID         string
	Name       string
	Email      string
	Password   string
	ProfileImg sql.NullString
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}
