package response

import (
	"courseworker/internal/dto"
	_error "courseworker/pkg/error"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return fieldName
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				return fieldName
			}
			if commaIdx := stringIndex(jsonTag, ','); commaIdx != -1 {
				return jsonTag[:commaIdx]
			}
			return jsonTag
		}
	}
	return fieldName
}

func stringIndex(s string, sep rune) int {
	for i, c := range s {
		if c == sep {
			return i
		}
	}
	return -1
}

func HttpBindingError(c *gin.Context, err error, reqStruct interface{}) {
	requestID := uuid.New().String()
	var logError LogError
	logError.Method = c.Request.Method
	logError.Path = c.Request.URL.Path
	logError.RequestID = requestID

	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		logError.Status = http.StatusBadRequest
		logError.Kind = _error.InvalidRequest.String()
		logError.Error = syntaxErr.Error()
		slog.With("data", logError).Error("Invalid request")

		c.JSON(http.StatusBadRequest, &ResponseError{
			Status:  false,
			Message: "Invalid request",
			Error: &ServiceError{
				RequestID: requestID,
				Kind:      _error.InvalidRequest.String(),
				Detail:    "the request payload contains invalid JSON format - please correct it",
			},
		})
		return
	} else if unmarshalTypeErr, ok := err.(*json.UnmarshalTypeError); ok {
		params := []_error.ProblemParameter{
			{
				Name:   getJSONFieldName(reflect.TypeOf(reqStruct), unmarshalTypeErr.Field),
				Reason: fmt.Sprintf("expected type '%s' but got '%s'", unmarshalTypeErr.Type, unmarshalTypeErr.Value),
			},
		}
		logError.Status = http.StatusBadRequest
		logError.Kind = _error.InvalidRequest.String()
		logError.Error = unmarshalTypeErr.Error()
		logError.Params = params
		slog.With("data", logError).Error("Invalid request")
		c.JSON(http.StatusBadRequest, &ResponseError{
			Status:  false,
			Message: "Invalid request",
			Error: &ServiceError{
				RequestID: requestID,
				Kind:      _error.InvalidRequest.String(),
				Detail:    "the request payload contains type mismatch",
				Param:     params,
			},
		})
		return
	} else if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var params []_error.ProblemParameter
		for _, fieldErr := range validationErrs {
			params = append(params, _error.ProblemParameter{
				Name:   getJSONFieldName(reflect.TypeOf(reqStruct), fieldErr.Field()),
				Reason: validationReasonMessage(fieldErr),
			})
		}
		logError.Status = http.StatusUnprocessableEntity
		logError.Kind = _error.Validation.String()
		logError.Error = validationErrs.Error()
		logError.Params = params
		slog.With("data", logError).Error("Invalid request")
		c.JSON(http.StatusUnprocessableEntity, &ResponseError{
			Status:  false,
			Message: "Invalid request",
			Error: &ServiceError{
				RequestID: requestID,
				Kind:      _error.Validation.String(),
				Detail:    "The request body contains failed field validation",
				Param:     params,
			},
		})
		return
	}
	logError.Status = http.StatusUnprocessableEntity
	logError.Kind = _error.InvalidRequest.String()
	logError.Error = err.Error()
	slog.With("data", logError).Error("Unexpected error")
	c.JSON(http.StatusInternalServerError, &ResponseError{
		Status:  false,
		Message: "Unexpected error",
		Error: &ServiceError{
			RequestID: requestID,
		},
	})
}

func validationReasonMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email format"
	case "url":
		return "must be a valid URL format"
	default:
		return fmt.Sprintf("failed validation for tag '%s'", fieldErr.Tag())
	}
}
