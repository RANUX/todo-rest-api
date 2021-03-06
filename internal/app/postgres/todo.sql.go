// Code generated by sqlc. DO NOT EDIT.
// source: todo.sql

package postgres

import (
	"context"
	"database/sql"
)

const createTodo = `-- name: CreateTodo :one
INSERT INTO "todos" (name, completed)
VALUES ($1, $2) RETURNING id, name, completed
`

type CreateTodoParams struct {
	Name      string       `json:"name"`
	Completed sql.NullBool `json:"completed"`
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.queryRow(ctx, q.createTodoStmt, createTodo, arg.Name, arg.Completed)
	var i Todo
	err := row.Scan(&i.ID, &i.Name, &i.Completed)
	return i, err
}

const deleteTodoById = `-- name: DeleteTodoById :exec
DELETE FROM todos WHERE id = $1
`

func (q *Queries) DeleteTodoById(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteTodoByIdStmt, deleteTodoById, id)
	return err
}

const getAllTodos = `-- name: GetAllTodos :many
SELECT id, name, completed FROM "todos" ORDER BY id
`

func (q *Queries) GetAllTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.query(ctx, q.getAllTodosStmt, getAllTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(&i.ID, &i.Name, &i.Completed); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTodoById = `-- name: GetTodoById :one
SELECT id, name, completed FROM "todos" WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTodoById(ctx context.Context, id int64) (Todo, error) {
	row := q.queryRow(ctx, q.getTodoByIdStmt, getTodoById, id)
	var i Todo
	err := row.Scan(&i.ID, &i.Name, &i.Completed)
	return i, err
}

const updateTodo = `-- name: UpdateTodo :one
UPDATE todos SET name = $2, completed = $3 WHERE id = $1 RETURNING id, name, completed
`

type UpdateTodoParams struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Completed sql.NullBool `json:"completed"`
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error) {
	row := q.queryRow(ctx, q.updateTodoStmt, updateTodo, arg.ID, arg.Name, arg.Completed)
	var i Todo
	err := row.Scan(&i.ID, &i.Name, &i.Completed)
	return i, err
}
