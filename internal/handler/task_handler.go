package handler

import (
	"courseworker/internal/dto"
	"courseworker/internal/service"
	_error "courseworker/pkg/error"
	"courseworker/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	serv service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{s}
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	resp, err := h.serv.GetAllTasksOfUser(claims.ID)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, tasksFetchSuccess, resp)
}

func (h *TaskHandler) GetTasksByCourse(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/GetTasksByCourse"), _error.InvalidRequest,
			_error.Title("Failed to get tasks"), "courseId must be a number",
		))
		return
	}

	resp, err := h.serv.GetTasksByCourseID(c, claims.ID, int64(courseID))
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, tasksFetchSuccess, resp)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	taskID := c.Param("taskId")
	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/GetTaskByID"), _error.InvalidRequest,
			_error.Title("Failed to get tasks"), "courseId must be a number",
		))
		return
	}

	resp, err := h.serv.GetTaskByID(c, claims.ID, taskID, int64(courseID))
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, taskFetchSuccess, resp)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/GetTaskByID"), _error.InvalidRequest,
			_error.Title("Failed to get tasks"), "courseId must be a number",
		))
		return
	}

	var req dto.TaskCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HttpBindingError(c, err, req)
	}

	resp, err := h.serv.CreateTask(c, claims.ID, int64(courseID), req)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, taskCreateSuccess, resp)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	taskID := c.Param("taskId")
	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/GetTaskByID"), _error.InvalidRequest,
			_error.Title("Failed to get tasks"), "courseId must be a number",
		))
		return
	}

	if err = h.serv.DeleteTask(c, claims.ID, taskID, int64(courseID)); err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, taskDeleteSuccess, nil)
}
