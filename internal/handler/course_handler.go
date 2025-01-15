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
		return
	}
	response.Success(c, http.StatusOK, "Course retrieved successfully", resp)
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	var req dto.CourseCreateUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HttpBindingError(c, err, req)
		return
	}

	resp, err := h.serv.CreateCourse(claims.ID, req)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "Course successfully created", resp)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/UpdateCourse"),
			_error.InvalidRequest,
			_error.Title("Failed to convert id from params"),
		))
		return
	}

	var req dto.CourseCreateUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HttpBindingError(c, err, req)
		return
	}

	resp, err := h.serv.UpdateCourse(claims.ID, int64(courseID), req)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Course successfully updated", resp)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	auth, _ := c.Get("user")
	claims := auth.(*dto.UserClaims)

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		response.HttpError(c, _error.E(
			_error.Op("hand/UpdateCourse"),
			_error.InvalidRequest,
			_error.Title("Failed to convert id from params"),
		))
		return
	}

	if err := h.serv.DeleteCourse(claims.ID, int64(courseID)); err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Course successfully deleted", nil)
}
