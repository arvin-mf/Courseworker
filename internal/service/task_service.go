package service

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	_error "courseworker/pkg/error"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TaskService interface {
	GetAllTasksOfUser(authUserID string) ([]dto.TaskResponse, error)
	GetTasksByCourseID(c *gin.Context, authUserID string, courseID int64) ([]dto.TaskResponse, error)
	GetTaskByID(c *gin.Context, authUserID, taskID string, courseID int64) (*dto.TaskResponse, error)
	CreateTask(c *gin.Context, authUserID string, courseID int64, req dto.TaskCreateReq) (*dto.ResponseID, error)
	DeleteTask(c *gin.Context, authUserID, taskID string, courseID int64) error
	SwitchTaskHighlight(c *gin.Context, authUserID, taskID string, courseID int64) (*dto.ResponseID, error)
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

func (s *taskService) CreateTask(c *gin.Context, authUserID string, courseID int64, req dto.TaskCreateReq) (*dto.ResponseID, error) {
	const op _error.Op = "serv/CreateTask"

	if err := s.cs.ValidateOwnershipCourse(c, authUserID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Forbidden action"), err)
	}

	deadline, err := time.Parse("2006-01-02 15:04", req.Deadline)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to create task"), err)
	}
	param := sqlc.CreateTaskParams{
		ID:          uuid.New().String(),
		CourseID:    courseID,
		Title:       req.Title,
		Type:        req.Type,
		Description: sql.NullString{String: req.Description, Valid: true},
		Deadline:    sql.NullTime{Time: deadline, Valid: true},
	}
	_, err = s.repo.CreateTask(param)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to create task"), err)
	}

	key := "task:" + authUserID
	if err = s.rd.Set(c, key, authUserID, 0).Err(); err != nil {
		log.Printf("Redis Set failed: %v", err)
	}

	return &dto.ResponseID{ID: param.ID}, nil
}

func (s *taskService) DeleteTask(c *gin.Context, authUserID, taskID string, courseID int64) error {
	const op _error.Op = "serv/DeleteTask"

	if err := s.ValidateOwnershipTask(c, authUserID, taskID, courseID); err != nil {
		return _error.E(op, _error.Forbidden, _error.Title("Failed to delete task"), err)
	}

	_, err := s.repo.DeleteTask(taskID)
	if err != nil {
		return _error.E(op, _error.Title("Failed to delete task"), err)
	}

	key := "task:" + taskID
	if err = s.rd.Del(c, key).Err(); err != nil {
		log.Printf("Redis Delete failed: %v", err)
	}

	return nil
}

func (s *taskService) SwitchTaskHighlight(c *gin.Context, authUserID, taskID string, courseID int64) (*dto.ResponseID, error) {
	const op _error.Op = "serv/SwitchTaskHighlight"

	if err := s.ValidateOwnershipTask(c, authUserID, taskID, courseID); err != nil {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Forbidden action"), err)
	}

	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get task"), err)
	}
	param := sqlc.SwitchTaskHighlightParams{
		Highlight: !task.Highlight,
		ID:        taskID,
	}
	_, err = s.repo.UpdateTaskHighlight(param)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to update task"), err)
	}
	return &dto.ResponseID{ID: param.ID}, nil
}
