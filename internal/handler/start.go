package handler

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/repository"
	"courseworker/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func route(r *gin.Engine, uh *UserHandler) {
	r.GET("/users", uh.GetUsers)
	r.GET("/users/:userId", uh.GetUserByID)
	r.GET("/auth/google/login-w-google", uh.LoginWithGoogle)
	r.GET("/auth/google/callback", uh.GetGoogleDetails)
	r.POST("/users", uh.RegisterUser)
	r.GET("/account-confirm", uh.CreateConfirmedUser)
}

func InitHandler(db *sql.DB) *UserHandler {
	queries := sqlc.New(db)

	userRepo := repository.NewUserRepository(queries)
	userServ := service.NewUserService(userRepo)
	userHand := NewUserHandler(userServ)

	return userHand
}

func StartEngine(r *gin.Engine, db *sql.DB) {
	uh := InitHandler(db)
	route(r, uh)
}
