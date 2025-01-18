// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: task.sql

package sqlc

import (
	"context"
	"database/sql"
)

const addImage = `-- name: AddImage :execresult
UPDATE tasks SET image = ? WHERE id = ?
`

type AddImageParams struct {
	Image sql.NullString
	ID    string
}

func (q *Queries) AddImage(ctx context.Context, arg AddImageParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addImage, arg.Image, arg.ID)
}

const createTask = `-- name: CreateTask :execresult
INSERT INTO tasks (id, course_id, title, type, description, deadline)
VALUES (?, ?, ?, ?, ?, ?)
`

type CreateTaskParams struct {
	ID          string
	CourseID    int64
	Title       string
	Type        string
	Description sql.NullString
	Deadline    sql.NullTime
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTask,
		arg.ID,
		arg.CourseID,
		arg.Title,
		arg.Type,
		arg.Description,
		arg.Deadline,
	)
}

const deleteTask = `-- name: DeleteTask :execresult
DELETE FROM tasks WHERE id = ?
`

func (q *Queries) DeleteTask(ctx context.Context, id string) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteTask, id)
}

const getAllTasks = `-- name: GetAllTasks :many
SELECT t.id, t.course_id, t.is_done, t.title, t.description, t.image, t.type, t.deadline, t.created_at, t.updated_at FROM tasks t
INNER JOIN courses c ON t.course_id = c.id
WHERE c.user_id = ?
`

func (q *Queries) GetAllTasks(ctx context.Context, userID string) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getAllTasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.CourseID,
			&i.IsDone,
			&i.Title,
			&i.Description,
			&i.Image,
			&i.Type,
			&i.Deadline,
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

const getTaskByID = `-- name: GetTaskByID :one
SELECT id, course_id, is_done, title, description, image, type, deadline, created_at, updated_at FROM tasks WHERE id = ?
`

func (q *Queries) GetTaskByID(ctx context.Context, id string) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.IsDone,
		&i.Title,
		&i.Description,
		&i.Image,
		&i.Type,
		&i.Deadline,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTasksByCourseID = `-- name: GetTasksByCourseID :many
SELECT id, course_id, is_done, title, description, image, type, deadline, created_at, updated_at FROM tasks WHERE course_id = ?
`

func (q *Queries) GetTasksByCourseID(ctx context.Context, courseID int64) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByCourseID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.CourseID,
			&i.IsDone,
			&i.Title,
			&i.Description,
			&i.Image,
			&i.Type,
			&i.Deadline,
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

const getUserIDFromTask = `-- name: GetUserIDFromTask :one
SELECT c.user_id FROM courses c
INNER JOIN tasks t ON t.course_id = c.id
WHERE t.id = ?
`

func (q *Queries) GetUserIDFromTask(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserIDFromTask, id)
	var user_id string
	err := row.Scan(&user_id)
	return user_id, err
}

const removeImage = `-- name: RemoveImage :execresult
UPDATE tasks SET image = NULL WHERE id = ?
`

func (q *Queries) RemoveImage(ctx context.Context, id string) (sql.Result, error) {
	return q.db.ExecContext(ctx, removeImage, id)
}

const updateTask = `-- name: UpdateTask :execresult
UPDATE tasks
SET title = ?, type = ?, description = ?, deadline = ?
WHERE id = ?
`

type UpdateTaskParams struct {
	Title       string
	Type        string
	Description sql.NullString
	Deadline    sql.NullTime
	ID          string
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTask,
		arg.Title,
		arg.Type,
		arg.Description,
		arg.Deadline,
		arg.ID,
	)
}
