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

type CourseHandler struct {
	serv service.CourseService
}

func NewCourseHandler(s service.CourseService) *CourseHandler {
	return &CourseHandler{s}
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	resp, err := h.serv.GetCoursesOfUser(claims.ID)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Course(s) retrieved successfully", resp)
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/GetCourseByID"),
			_error.InvalidRequest,
			_error.Title("Failed to get course"),
			_error.Detail("failed parsing course id"),
		))
		return
	}

	resp, err := h.serv.GetCourseByID(claims.ID, int64(courseID))
	if err != nil {
		response.HttpError(c, err)
	}
	response.Success(c, http.StatusOK, "Course retrieved successfully", resp)
}
