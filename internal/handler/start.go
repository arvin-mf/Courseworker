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

func route(r *gin.Engine, uh *UserHandler, ch *CourseHandler, th *TaskHandler) {
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

	r.GET("/courses/tasks", middleware.ValidateToken(), th.GetAllTasks)
	r.GET("/courses/:courseId/tasks", middleware.ValidateToken(), th.GetTasksByCourse)
	r.GET("/courses/:courseId/tasks/:taskId", middleware.ValidateToken(), th.GetTaskByID)
	r.POST("/courses/:courseId/tasks", middleware.ValidateToken(), th.CreateTask)
	r.PUT("/courses/:courseId/tasks/:taskId/highlight", middleware.ValidateToken(), th.SwitchTaskHighlight)
	r.DELETE("/courses/:courseId/tasks/:taskId", middleware.ValidateToken(), th.DeleteTask)
}

func InitHandler(db *sql.DB, rd *redis.Client) (*UserHandler, *CourseHandler, *TaskHandler) {
	queries := sqlc.New(db)

	userRepo := repository.NewUserRepository(queries)
	userServ := service.NewUserService(userRepo)
	userHand := NewUserHandler(userServ)

	courseRepo := repository.NewCourseRepository(queries)
	courseServ := service.NewCourseService(courseRepo, rd)
	courseHand := NewCourseHandler(courseServ)

	taskRepo := repository.NewTaskRepository(queries)
	taskServ := service.NewTaskService(taskRepo, rd, courseServ)
	taskHand := NewTaskHandler(taskServ)

	return userHand, courseHand, taskHand
}

func StartEngine(r *gin.Engine, db *sql.DB, rd *redis.Client) {
	uh, ch, th := InitHandler(db, rd)
	route(r, uh, ch, th)
}
