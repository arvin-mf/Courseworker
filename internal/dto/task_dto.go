package dto

import (
	"courseworker/internal/db/sqlc"
	"time"
)

type TaskResponse struct {
	ID          string    `json:"id"`
	CourseID    int64     `json:"course_id"`
	IsDone      bool      `json:"is_done"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Type        string    `json:"type"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToTaskResponse(t *sqlc.Task) *TaskResponse {
	return &TaskResponse{
		ID: t.ID, CourseID: t.CourseID, IsDone: t.IsDone,
		Title: t.Title, Description: t.Description.String,
		Image: t.Image.String, Type: t.Type, Deadline: t.Deadline.Time,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func ToTaskResponses(tasks *[]sqlc.Task) []TaskResponse {
	responses := []TaskResponse{}
	for _, t := range *tasks {
		response := TaskResponse{
			ID: t.ID, CourseID: t.CourseID, IsDone: t.IsDone,
			Title: t.Title, Description: t.Description.String,
			Image: t.Image.String, Type: t.Type, Deadline: t.Deadline.Time,
			CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return responses
}

type TaskCreateReq struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}
