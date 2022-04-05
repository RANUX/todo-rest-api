package apiserver

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ranux/todo-rest-api/internal/app/postgres"
)

type Handlers struct {
	Repo *postgres.Repo
}

func NewHandlers(repo *postgres.Repo) *Handlers {
	return &Handlers{Repo: repo}
}

func mapTodo(todo postgres.Todo) interface{} {
	return struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}{
		ID:        todo.ID,
		Name:      todo.Name,
		Completed: todo.Completed.Bool,
	}
}

func (h *Handlers) UpdateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return nil
	}

	todo, err := h.Repo.GetTodoById(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return nil
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = sql.NullBool{
			Bool:  *body.Completed,
			Valid: true,
		}
	}

	todo, err = h.Repo.UpdateTodo(ctx.Context(), postgres.UpdateTodoParams{
		ID:        int64(id),
		Name:      todo.Name,
		Completed: todo.Completed,
	})

	if err != nil {
		ctx.SendStatus(fiber.StatusNotFound)
		return nil
	}

	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}
	return nil
}

func (h *Handlers) DeleteTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	_, err = h.Repo.GetTodoById(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return nil
	}

	err = h.Repo.DeleteTodoById(ctx.Context(), int64(id))
	if err != nil {
		ctx.SendStatus(fiber.StatusNotFound)
		return nil
	}

	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}

func (h *Handlers) GetTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	todo, err := h.Repo.GetTodoById(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return nil
	}

	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}
	return nil
}

func (h *Handlers) CreateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return nil
	}

	if len(body.Name) <= 2 {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is not long enough",
		})
		return nil
	}

	todo, err := h.Repo.CreateTodo(ctx.Context(), postgres.CreateTodoParams{
		Name: body.Name,
		Completed: sql.NullBool{
			Bool:  body.Completed,
			Valid: true,
		},
	})

	if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(mapTodo(todo)); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}
	return nil
}

func (h *Handlers) GetTodos(ctx *fiber.Ctx) error {
	todos, err := h.Repo.GetAllTodos(ctx.Context())
	if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}

	result := make([]interface{}, len(todos))
	for i, todo := range todos {
		result[i] = mapTodo(todo)
	}

	if err := ctx.Status(fiber.StatusOK).JSON(result); err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return nil
	}
	return nil
}
