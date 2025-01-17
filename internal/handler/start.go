package handler

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/repository"
	"courseworker/internal/service"
	"courseworker/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func route(r *gin.Engine, uh *UserHandler, ch *CourseHandler) {
	r.GET("/users", uh.GetUsers)
	r.GET("/users/:userId", uh.GetUserByID)
	r.GET("/auth/google/login-w-google", uh.LoginWithGoogle)
	r.GET("/auth/google/callback", uh.GetGoogleDetails)
	r.POST("/register", uh.RegisterUser)
	r.GET("/account-confirm", uh.CreateConfirmedUser)
	r.POST("/login", uh.LoginUser)

	r.GET("/courses", middleware.ValidateToken(), ch.GetCourses)
	r.GET("/courses/:courseId", middleware.ValidateToken(), ch.GetCourseByID)
	r.POST("/courses", middleware.ValidateToken(), ch.CreateCourse)
	r.PUT("/courses/:courseId", middleware.ValidateToken(), ch.UpdateCourse)
	r.DELETE("/courses/:courseId", middleware.ValidateToken(), ch.DeleteCourse)
}

func InitHandler(db *sql.DB, rd *redis.Client) (*UserHandler, *CourseHandler) {
	queries := sqlc.New(db)

	userRepo := repository.NewUserRepository(queries)
	userServ := service.NewUserService(userRepo)
	userHand := NewUserHandler(userServ)

	courseRepo := repository.NewCourseRepository(queries)
	courseServ := service.NewCourseService(courseRepo, rd)
	courseHand := NewCourseHandler(courseServ)

	return userHand, courseHand
}

func StartEngine(r *gin.Engine, db *sql.DB, rd *redis.Client) {
	uh, ch := InitHandler(db, rd)
	route(r, uh, ch)
}
