package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
	"database/sql"
	"errors"
)

type UserRepository interface {
	GetAllUsers() ([]sqlc.GetAllUsersRow, error)
	GetUserByID(userID string) (*sqlc.User, error)
}

type userRepository struct {
	db *sqlc.Queries
}

func NewUserRepository(db *sqlc.Queries) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAllUsers() ([]sqlc.GetAllUsersRow, error) {
	const op _error.Op = "repo/GetAllUsers"
	result, err := r.db.GetAllUsers(context.Background())
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *userRepository) GetUserByID(userID string) (*sqlc.User, error) {
	const op _error.Op = "repo/GetUserByID"
	result, err := r.db.GetUserByID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, _error.E(op, _error.NotExist, err)
		}
		return nil, _error.E(op, _error.Database, err)
	}
	return &result, nil
}
