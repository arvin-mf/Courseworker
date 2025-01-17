package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
	"database/sql"
	"errors"
	"fmt"
)

type CourseRepository interface {
	GetAllCourses(userID string) ([]sqlc.Course, error)
	GetCourseByID(ID int64) (*sqlc.Course, error)
	CreateCourse(param sqlc.CreateCourseParams) (sql.Result, error)
	UpdateCourse(param sqlc.UpdateCourseParams) (sql.Result, error)
	DeleteCourse(courseID int64) (sql.Result, error)
	GetUserIDFromCourse(courseID int64) (string, error)
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
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.Course{}, nil
		}
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

func (r *courseRepository) CreateCourse(param sqlc.CreateCourseParams) (sql.Result, error) {
	const op _error.Op = "repo/CreateCourse"
	result, err := r.db.CreateCourse(context.Background(), param)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *courseRepository) UpdateCourse(param sqlc.UpdateCourseParams) (sql.Result, error) {
	const op _error.Op = "repo/UpdateCourse"
	result, err := r.db.UpdateCourse(context.Background(), param)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *courseRepository) DeleteCourse(courseID int64) (sql.Result, error) {
	const op _error.Op = "repo/DeleteCourse"
	result, err := r.db.DeleteCourse(context.Background(), courseID)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	if affected == 0 {
		return nil, _error.E(
			op, _error.Title("No row affected"),
			fmt.Sprintf("The requested course with id %d could not be found", courseID),
		)
	}
	return result, nil
}

func (r *courseRepository) GetUserIDFromCourse(courseID int64) (string, error) {
	const op _error.Op = "repo/GetUserIDFromCourse"
	result, err := r.db.GetUserIDFromCourse(context.Background(), courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", _error.E(
				op, _error.NotExist, _error.Title("Course not found"),
				fmt.Sprintf("The requested course with id %d could not be found", courseID),
			)
		}
		return "", _error.E(op, _error.Database, err)
	}
	return result, nil
}
