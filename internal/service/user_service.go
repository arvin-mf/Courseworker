package service

import (
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	_error "courseworker/pkg/error"
)

type UserService interface {
	GetUsers() ([]dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetUsers() ([]dto.UserResponse, error) {
	const op _error.Op = "serv/GetUsers"
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get users"), err)
	}
	return dto.ToUserResponses(&users), nil
}
