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
	}
	response.Success(c, http.StatusOK, "Users retrieved successfully", resp)
}
