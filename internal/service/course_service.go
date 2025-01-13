package service

import (
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	_error "courseworker/pkg/error"
)

type CourseService interface {
	GetCoursesOfUser(userID string) ([]dto.CourseResponse, error)
	GetCourseByID(userID string, courseID int64) (*dto.CourseResponse, error)
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(r repository.CourseRepository) CourseService {
	return &courseService{r}
}

func (s *courseService) GetCoursesOfUser(userID string) ([]dto.CourseResponse, error) {
	const op _error.Op = "serv/GetCoursesOfUser"
	courses, err := s.repo.GetAllCourses(userID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get courses"), err)
	}
	return dto.ToCourseResponses(&courses), nil
}

func (s *courseService) GetCourseByID(userID string, courseID int64) (*dto.CourseResponse, error) {
	const op _error.Op = "serv/GetCourseByID"
	course, err := s.repo.GetCourseByID(courseID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get course"), err)
	}

	if userID != course.UserID {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Course does not belong"), err)
	}

	return dto.ToCourseResponse(course), nil
}
