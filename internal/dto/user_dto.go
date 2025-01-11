package dto

import (
	"courseworker/internal/db/sqlc"
	"time"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ProfileImg string    `json:"profile_img"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToUserResponses(users *[]sqlc.GetAllUsersRow) []UserResponse {
	var responses []UserResponse
	for _, u := range *users {
		response := UserResponse{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			ProfileImg: u.ProfileImg.String,
			CreatedAt:  u.CreatedAt.Time,
			UpdatedAt:  u.UpdatedAt.Time,
		}
		responses = append(responses, response)
	}
	return responses
}
