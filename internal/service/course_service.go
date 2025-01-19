package service

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	_error "courseworker/pkg/error"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CourseService interface {
	GetCoursesOfUser(userID string) ([]dto.CourseResponse, error)
	GetCourseByID(c *gin.Context, userID string, courseID int64) (*dto.CourseResponse, error)
	CreateCourse(c *gin.Context, userID string, arg dto.CourseCreateUpdateReq) (*dto.ResponseID, error)
	UpdateCourse(c *gin.Context, userID string, courseID int64, arg dto.CourseCreateUpdateReq) (*dto.ResponseID, error)
	DeleteCourse(c *gin.Context, userID string, courseID int64) error
	ValidateOwnershipCourse(c *gin.Context, authUserID string, courseID int64) error
}

type courseService struct {
	repo repository.CourseRepository
	rd   *redis.Client
}

func NewCourseService(r repository.CourseRepository, rdc *redis.Client) CourseService {
	return &courseService{
		repo: r,
		rd:   rdc,
	}
}

func (s *courseService) GetCoursesOfUser(userID string) ([]dto.CourseResponse, error) {
	const op _error.Op = "serv/GetCoursesOfUser"
	courses, err := s.repo.GetAllCourses(userID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get courses"), err)
	}
	return dto.ToCourseResponses(&courses), nil
}

func (s *courseService) ValidateOwnershipCourse(c *gin.Context, authUserID string, courseID int64) error {
	const op _error.Op = "serv/validateOwnershipCourse"

	var value string
	key := "course:" + strconv.Itoa(int(courseID))
	value, err := s.rd.Get(c, key).Result()
	if err != nil {
		log.Printf("Redis Get failed: %v", err)
		value, err = s.repo.GetUserIDFromCourse(courseID)
		if err != nil {
			return _error.E(op, _error.Title("Failed to get userID"), err)
		}
		if err = s.rd.Set(c, key, value, 0).Err(); err != nil {
			log.Printf("Redis Set failed: %v", err)
		}
	}

	if authUserID != value {
		return _error.E(
			op, _error.Forbidden, _error.Title("Forbidden action"),
			fmt.Sprintf("The requested course with id %d does not belong to user", courseID),
		)
	}
	return nil
}

func (s *courseService) GetCourseByID(c *gin.Context, userID string, courseID int64) (*dto.CourseResponse, error) {
	const op _error.Op = "serv/GetCourseByID"

	if err := s.ValidateOwnershipCourse(c, userID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Forbidden action"), err)
	}

	course, err := s.repo.GetCourseByID(courseID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get course"), err)
	}

	if userID != course.UserID {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Course does not belong"), err)
	}

	return dto.ToCourseResponse(course), nil
}

func (s *courseService) CreateCourse(c *gin.Context, userID string, arg dto.CourseCreateUpdateReq) (*dto.ResponseID, error) {
	const op _error.Op = "serv/CreateCourse"
	result, err := s.repo.CreateCourse(sqlc.CreateCourseParams{
		Name:    arg.Name,
		Subname: sql.NullString{String: arg.Subname, Valid: true},
		UserID:  userID,
	})
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to create course"), err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get new id"), err)
	}

	ctx := c.Request.Context()
	key := "course:" + strconv.Itoa(int(id))
	if err = s.rd.Set(ctx, key, userID, 0).Err(); err != nil {
		log.Printf("Redis Set failed: %v", err)
	}

	return &dto.ResponseID{ID: id}, nil
}

func (s *courseService) UpdateCourse(c *gin.Context, userID string, courseID int64, arg dto.CourseCreateUpdateReq) (*dto.ResponseID, error) {
	const op _error.Op = "serv/UpdateCourse"

	if err := s.ValidateOwnershipCourse(c, userID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Forbidden action"), err)
	}

	_, err := s.repo.UpdateCourse(sqlc.UpdateCourseParams{
		Name:    arg.Name,
		Subname: sql.NullString{String: arg.Subname, Valid: true},
		ID:      courseID,
	})
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to update course"), err)
	}
	return &dto.ResponseID{ID: courseID}, nil
}

func (s *courseService) DeleteCourse(c *gin.Context, userID string, courseID int64) error {
	const op _error.Op = "serv/DeleteCourse"

	if err := s.ValidateOwnershipCourse(c, userID, courseID); err != nil {
		return _error.E(op, _error.Forbidden, _error.Title("Forbidden action"), err)
	}

	_, err := s.repo.DeleteCourse(courseID)
	if err != nil {
		return _error.E(op, _error.Title("Failed to delete course"), err)
	}

	key := "course:" + strconv.Itoa(int(courseID))
	if err := s.rd.Del(c, key).Err(); err != nil {
		log.Printf("Redis Delete failed: %v", err)
	}

	return nil
}
