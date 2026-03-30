-- name: AddUser :execresult
INSERT INTO users (name,email, password)
VALUES(?,?,?);

-- name: GetUserById :one
SELECT id, name, email, created_at
FROM users WHERE id = ?;


-- name: CreateTask :execresult
INSERT INTO tasks (title, description, status_id, user_id, start_time, end_time)
VALUES(?,?,?,?,?,?);

-- name: GetTask :one
SELECT 
    t.id, 
    t.title, 
    t.description, 
    s.status AS status_name,
    t.user_id,
    t.start_time,
    t.end_time,
    t.created_at, 
    t.updated_at
FROM tasks t
JOIN master_status s ON t.status_id = s.id
WHERE t.id = ?;

-- name: ListTasks :many
SELECT 
    t.id, sql/queries.sql:27:7: column reference "id" is ambiguous
    t.title, 
    t.description, 
    s.status AS status_name,
    t.start_time,
    t.end_time,
    t.user_id,
    t.created_at, 
    t.updated_at
FROM tasks t
JOIN master_status s ON t.status_id = s.id;

-- name: UpdateTasks :exec 
UPDATE tasks
SET title = ?, description = ?, status_id = ?, start_time = ?, end_time = ?
WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

