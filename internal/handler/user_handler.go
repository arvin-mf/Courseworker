package handler

import (
	"courseworker/internal/service"
	"courseworker/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	serv service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	resp, err := h.serv.GetUsers()
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Users retrieved successfully", resp)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("userId")
	resp, err := h.serv.GetUserByID(userID)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "User retrieved successfully", resp)
}
