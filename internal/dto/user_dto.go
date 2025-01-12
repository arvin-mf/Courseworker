package dto

import (
	"courseworker/internal/db/sqlc"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ProfileImg string    `json:"profile_img"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToUserResponse(u *sqlc.User) *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		ProfileImg: u.ProfileImg.String,
		CreatedAt:  u.CreatedAt.Time,
		UpdatedAt:  u.UpdatedAt.Time,
	}
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

type UserClaims struct {
	ID          string `json:"id" binding:"required"`
	ExpDuration int64  `json:"exp_duration"`
	jwt.RegisteredClaims
}

func NewUserClaims(ID string, exp time.Duration) UserClaims {
	return UserClaims{
		ID:          ID,
		ExpDuration: int64(exp.Seconds()),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}

type AuthenticatedUser struct {
	ID            string `json:"localId"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Name          string `json:"name"`
	Photo         string `json:"photoUrl"`
}

type RegisterUserReq struct {
	Name            string `json:"name"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type RegisterUserResp struct {
	Email string `json:"email"`
}

type RegistrationClaims struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	HashedPw string `json:"hashed_pw" binding:"required"`
	jwt.RegisteredClaims
}

func NewRegistrationClaims(name, email, pass string, exp time.Duration) RegistrationClaims {
	return RegistrationClaims{
		Name:     name,
		Email:    email,
		HashedPw: pass,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
