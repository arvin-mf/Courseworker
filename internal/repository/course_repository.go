package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
	"database/sql"
	"errors"
)

type CourseRepository interface {
	GetAllCourses(userID string) ([]sqlc.Course, error)
	GetCourseByID(ID int64) (*sqlc.Course, error)
}

type courseRepository struct {
	db *sqlc.Queries
}

func NewCourseRepository(db *sqlc.Queries) CourseRepository {
	return &courseRepository{db}
}

func (r *courseRepository) GetAllCourses(userID string) ([]sqlc.Course, error) {
	const op _error.Op = "repo/GetAllCourses"
	result, err := r.db.GetAllCourses(context.Background(), userID)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *courseRepository) GetCourseByID(ID int64) (*sqlc.Course, error) {
	const op _error.Op = "repo/GetCourseByID"
	result, err := r.db.GetCourseByID(context.Background(), ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, _error.E(op, _error.NotExist, err)
		}
		return nil, _error.E(op, _error.Database, err)
	}
	return &result, nil
}
