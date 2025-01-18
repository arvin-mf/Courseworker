-- name: GetAllTasks :many
SELECT t.* FROM tasks t
INNER JOIN courses c ON t.course_id = c.id
WHERE c.user_id = ?;

-- name: GetTasksByCourseID :many
SELECT * FROM tasks WHERE course_id = ?;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = ?;

-- name: GetUserIDFromTask :one
SELECT c.user_id FROM courses c
INNER JOIN tasks t ON t.course_id = c.id
WHERE t.id = ?;

-- name: CreateTask :execresult
INSERT INTO tasks (id, course_id, title, type, description, deadline)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateTask :execresult
UPDATE tasks
SET title = ?, type = ?, description = ?, deadline = ?
WHERE id = ?;

-- name: AddImage :execresult
UPDATE tasks SET image = ? WHERE id = ?;

-- name: RemoveImage :execresult
UPDATE tasks SET image = NULL WHERE id = ?;

-- name: DeleteTask :execresult
DELETE FROM tasks WHERE id = ?;
