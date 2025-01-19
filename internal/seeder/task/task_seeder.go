package main

import (
	"context"
	"courseworker/config"
	"courseworker/internal/db/sqlc"
	"courseworker/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.SetupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rdc := config.NewRedisClient()
	defer rdc.Close()

	queries := sqlc.New(db)
	repo_t := repository.NewTaskRepository(queries)
	repo_c := repository.NewCourseRepository(queries)
	repo_u := repository.NewUserRepository(queries)

	fmt.Println("User's email: ")
	var email string
	fmt.Scan(&email)

	user, err := repo_u.GetUserByEmail(email)
	if err != nil {
		log.Fatal("Failed to get user")
	}

	courseIDs := []int{}
	for i := 0; i < 4; i++ {
		result, err := repo_c.CreateCourse(sqlc.CreateCourseParams{
			Name:    fmt.Sprintf("Course %d", i+1),
			Subname: sql.NullString{String: "Lorem ipsum dolor sit amet", Valid: true},
			UserID:  user.ID,
		})
		if err != nil {
			log.Fatal("Failed to create course")
		}

		courseID, err := result.LastInsertId()
		if err != nil {
			log.Fatal("Failed to get course id")
		}
		key := "course:" + strconv.Itoa(int(courseID))
		if err = rdc.Set(context.Background(), key, user.ID, 0).Err(); err != nil {
			log.Printf("Redis Set failed: %v", err)
		}

		courseIDs = append(courseIDs, int(courseID))
	}

	for _, c := range courseIDs {
		for i := 0; i < 3; i++ {
			taskParam := sqlc.CreateTaskParams{
				ID:       uuid.New().String(),
				CourseID: int64(c),
				Title:    fmt.Sprintf("Task %d of Course%d", i+1, c),
				Type:     "Individual",
				Description: sql.NullString{
					String: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin in tellus ac dui suscipit imperdiet. Aliquam bibendum ipsum mi, vel feugiat nunc lacinia a.",
					Valid:  true,
				},
				Deadline: sql.NullTime{Time: time.Now().Add(120 * time.Hour), Valid: true},
			}
			_, err := repo_t.CreateTask(taskParam)
			if err != nil {
				log.Fatal("Failed to create task")
			}

			key := "task:" + taskParam.ID
			if err = rdc.Set(context.Background(), key, user.ID, 0).Err(); err != nil {
				log.Printf("Redis Set failed: %v", err)
			}
		}
	}

	fmt.Println("Database successfully seeded")
}
