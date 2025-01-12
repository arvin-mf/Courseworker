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
	EmailExists(email string) (int64, error)
	GetUserByEmail(email string) (*sqlc.User, error)
	CreateUser(sqlc.CreateUserParams) (sql.Result, error)
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

func (r *userRepository) EmailExists(email string) (int64, error) {
	const op _error.Op = "repo/EmailExists"
	result, err := r.db.CountUserByEmail(context.Background(), email)
	if err != nil {
		return -1, _error.E(op, _error.Database, err)
	}
	return result, nil
}

func (r *userRepository) GetUserByEmail(email string) (*sqlc.User, error) {
	const op _error.Op = "repo/GetUserByEmail"
	result, err := r.db.GetUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, _error.E(op, _error.NotExist, err)
		}
		return nil, _error.E(op, _error.Database, err)
	}
	return &result, nil
}

func (r *userRepository) CreateUser(param sqlc.CreateUserParams) (sql.Result, error) {
	const op _error.Op = "repo/CreateUser"
	result, err := r.db.CreateUser(context.Background(), param)
	if err != nil {
		return nil, _error.E(op, _error.Database, err)
	}
	return result, nil
}
