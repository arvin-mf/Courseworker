// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package sqlc

import (
	"database/sql"
	"time"
)

type Course struct {
	ID        int64
	Name      string
	Subname   sql.NullString
	UserID    string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type Task struct {
	ID          string
	CourseID    int64
	IsDone      bool
	Title       string
	Description sql.NullString
	Image       sql.NullString
	Type        string
	Deadline    sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskNote struct {
	ID        int64
	TaskID    string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID         string
	Name       string
	Email      string
	Password   string
	ProfileImg sql.NullString
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}
