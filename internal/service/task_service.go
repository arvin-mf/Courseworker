package service

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	_error "courseworker/pkg/error"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TaskService interface {
	GetAllTasksOfUser(authUserID string) ([]dto.TaskResponse, error)
	GetTasksByCourseID(c *gin.Context, authUserID string, courseID int64) ([]dto.TaskResponse, error)
	GetTaskByID(c *gin.Context, authUserID, taskID string, courseID int64) (*dto.TaskResponse, error)
}

type taskService struct {
	repo repository.TaskRepository
	rd   *redis.Client
	cs   CourseService
}

func NewTaskService(r repository.TaskRepository, rdc *redis.Client, courseServ CourseService) TaskService {
	return &taskService{
		repo: r,
		rd:   rdc,
		cs:   courseServ,
	}
}

func (s *taskService) GetAllTasksOfUser(authUserID string) ([]dto.TaskResponse, error) {
	const op _error.Op = "serv/GetAllTasksOfUser"
	tasks, err := s.repo.GetAllTasks(authUserID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get tasks"), err)
	}
	return dto.ToTaskResponses(&tasks), nil
}

func (s *taskService) GetTasksByCourseID(c *gin.Context, authUserID string, courseID int64) ([]dto.TaskResponse, error) {
	const op _error.Op = "serv/GetTasksByCourseID"

	if err := s.cs.ValidateOwnershipCourse(c, authUserID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Failed to get tasks"), err)
	}

	tasks, err := s.repo.GetTasksByCourse(courseID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get tasks"), err)
	}
	return dto.ToTaskResponses(&tasks), nil
}

func (s *taskService) GetTaskByID(c *gin.Context, authUserID, taskID string, courseID int64) (*dto.TaskResponse, error) {
	const op _error.Op = "serv/GetTaskByID"

	if err := s.ValidateOwnershipTask(c, authUserID, taskID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Failed to get task"), err)
	}

	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get task"), err)
	}
	return dto.ToTaskResponse(task), nil
}

func (s *taskService) ValidateOwnershipTask(c *gin.Context, authUserID, taskID string, courseID int64) error {
	const op _error.Op = "serv/validateOwnershipTask"

	var value string
	key := "task:" + taskID
	value, err := s.rd.Get(c, key).Result()
	if err != nil {
		log.Printf("Redis Get failed: %v", err)
		value, err = s.repo.GetUserIDFromTask(sqlc.GetUserIDFromTaskParams{
			TaskID:   taskID,
			CourseID: courseID,
		})
		if err != nil {
			return _error.E(op, _error.Cache, _error.Title("Failed to get userID"), err)
		}
		if err = s.rd.Set(c, key, value, 0).Err(); err != nil {
			log.Printf("Redis Set failed: %v", err)
		}
	}

	if authUserID != value {
		return _error.E(
			op, _error.Forbidden, _error.Title("Forbidden action"),
			fmt.Sprintf("The requested task with id %s does not belong to user", taskID),
		)
	}
	return nil
}
