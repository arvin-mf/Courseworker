package response

import (
	"courseworker/internal/dto"
	_error "courseworker/pkg/error"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Success(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, dto.Response{
		Status:  true,
		Message: msg,
		Data:    data,
	})
}

type ServiceError struct {
	RequestID string      `json:"request_id,omitempty"`
	Kind      string      `json:"kind,omitempty"`
	Detail    string      `json:"detail,omitempty"`
	Param     interface{} `json:"param,omitempty"`
}

type ResponseError struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Error   *ServiceError `json:"error,omitempty"`
}

func getKindStatusCode(arg _error.Kind) int {
	switch arg {
	case _error.NotExist:
		return http.StatusNotFound
	case _error.Forbidden:
		return http.StatusForbidden
	case _error.InvalidRequest:
		return http.StatusBadRequest
	case _error.Database, _error.Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

type LogError struct {
	RequestID string      `json:"request_id,omitempty"`
	Status    int         `json:"status,omitempty"`
	Method    string      `json:"method,omitempty"`
	Path      string      `json:"path,omitempty"`
	Stack     []string    `json:"stack,omitempty"`
	Kind      string      `json:"kind,omitempty"`
	Params    interface{} `json:"params,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func HttpError(c *gin.Context, err error) {
	requestID := uuid.New().String()
	var problem *_error.Problem
	if errors.As(err, &problem) {
		status := getKindStatusCode(problem.Kind)
		opStack := _error.OpStack(problem)

		logData := &LogError{
			RequestID: requestID,
			Status:    status,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Stack:     opStack,
			Kind:      problem.Kind.String(),
			Error:     problem.Error(),
		}

		slog.With("data", logData).Error(string(problem.Title))

		if status == 500 {
			c.JSON(status, &ResponseError{
				Status:  false,
				Message: string(problem.Title),
				Error: &ServiceError{
					RequestID: requestID,
					Kind:      problem.Kind.String(),
					Detail:    "Internal server error has occured - please contact support",
				},
			})
			return
		} else {
			detail := string(problem.Detail)
			if detail == "" {
				detail = problem.Error()
			}
			c.JSON(status, &ResponseError{
				Status:  false,
				Message: string(problem.Title),
				Error: &ServiceError{
					RequestID: requestID,
					Kind:      problem.Kind.String(),
					Detail:    detail,
				},
			})
			return
		}
	}
	logData := &LogError{
		RequestID: requestID,
		Status:    http.StatusInternalServerError,
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Stack:     []string{"unexpected error occured"},
		Kind:      _error.Other.String(),
		Error:     err.Error(),
	}

	slog.With("data", logData).Error(string(problem.Title))
	c.JSON(http.StatusInternalServerError, &ResponseError{
		Status:  false,
		Message: "Unexpected error",
		Error: &ServiceError{
			RequestID: requestID,
		},
	})
}
