-- name: GetAllTodos :many
SELECT * FROM "todos";

-- name: CreateTodo :one
INSERT INTO "todos" (name, completed)
VALUES ($1, $2) RETURNING *;

-- name: GetTodoById :one
SELECT * FROM "todos" WHERE id = $1 LIMIT 1;

