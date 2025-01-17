-- name: GetAllCourses :many
SELECT * FROM courses WHERE user_id = ?;

-- name: GetCourseByID :one
SELECT * FROM courses WHERE id = ?;

-- name: CreateCourse :execresult
INSERT INTO courses (name, subname, user_id)
VALUES (?, ?, ?);

-- name: UpdateCourse :execresult
UPDATE courses
SET name = ?, subname = ?
WHERE id = ?;

-- name: DeleteCourse :execresult
DELETE FROM courses WHERE id = ?;

-- name: GetUserIDFromCourse :one
SELECT user_id FROM courses WHERE id = ?;