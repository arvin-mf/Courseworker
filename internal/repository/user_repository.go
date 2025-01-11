package repository

import (
	"context"
	"courseworker/internal/db/sqlc"
	_error "courseworker/pkg/error"
)

type UserRepository interface {
	GetAllUsers() ([]sqlc.GetAllUsersRow, error)
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
