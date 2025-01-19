package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
	"database/sql"
	"errors"
)

type TaskRepository interface {
	GetAllTasks(userID string) ([]sqlc.Task, error)
	GetTasksByCourse(courseID int64) ([]sqlc.Task, error)
	GetTaskByID(taskID string) (*sqlc.Task, error)
}

type taskRepository struct {
	db *sqlc.Queries
}

func NewTaskRepository(db *sqlc.Queries) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetAllTasks(userID string) ([]sqlc.Task, error) {
	const op _error.Op = "repo/GetAllTasks"
	result, err := r.db.GetAllTasks(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.Task{}, nil
		}
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *taskRepository) GetTasksByCourse(courseID int64) ([]sqlc.Task, error) {
	const op _error.Op = "repo/GetTasksByCourse"
	result, err := r.db.GetTasksByCourseID(context.Background(), courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.Task{}, nil
		}
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *taskRepository) GetTaskByID(taskID string) (*sqlc.Task, error) {
	const op _error.Op = "repo/GetTaskByID"
	result, err := r.db.GetTaskByID(context.Background(), taskID)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return &result, nil
}
