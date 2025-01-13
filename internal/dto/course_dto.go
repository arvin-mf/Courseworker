package dto

import (
	"courseworker/internal/db/sqlc"
	"time"
)

type CourseResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Subname   string    `json:"subname"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCourseResponse(c *sqlc.Course) *CourseResponse {
	return &CourseResponse{
		ID:        c.ID,
		Name:      c.Name,
		Subname:   c.Subname.String,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}

func ToCourseResponses(courses *[]sqlc.Course) []CourseResponse {
	var responses []CourseResponse
	for _, c := range *courses {
		response := CourseResponse{
			ID:        c.ID,
			Name:      c.Name,
			Subname:   c.Subname.String,
			UserID:    c.UserID,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		}
		responses = append(responses, response)
	}
	return responses
}
