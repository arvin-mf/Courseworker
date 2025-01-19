package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
	"database/sql"
	"errors"
	"fmt"
)

type TaskRepository interface {
	GetAllTasks(userID string) ([]sqlc.Task, error)
	GetTasksByCourse(courseID int64) ([]sqlc.Task, error)
	GetTaskByID(taskID string) (*sqlc.Task, error)
	GetUserIDFromTask(param sqlc.GetUserIDFromTaskParams) (string, error)
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

func (r *taskRepository) GetUserIDFromTask(param sqlc.GetUserIDFromTaskParams) (string, error) {
	const op _error.Op = "repo/GetUserIDFromTask"
	result, err := r.db.GetUserIDFromTask(context.Background(), param)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", _error.E(
				op, _error.NotExist, _error.Title("Task not found"),
				fmt.Sprintf("The requested task with id %s could not be found", param.TaskID),
			)
		}
		return "", nil
	}
	return result, nil
}
