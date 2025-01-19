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

}
