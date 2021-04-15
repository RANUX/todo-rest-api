package apiserver

import (
	"strconv"

	"github.com/gofiber/fiber"
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

// func (h *Handlers) UpdateTodo(ctx *fiber.Ctx) {
// 	type request struct {
// 		Name      *string `json:"name"`
// 		Completed *bool   `json:"completed"`
// 	}

// 	paramsId := ctx.Params("id")
// 	id, err := strconv.Atoi(paramsId)
// 	if err != nil {
// 		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "cannot parse id",
// 		})
// 		return
// 	}

// 	var body request
// 	err = ctx.BodyParser(&body)
// 	if err != nil {
// 		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "cannot parse body",
// 		})
// 		return
// 	}

// 	todo, err := h.Repo.GetTodoById(ctx.Context(), int64(id))
// 	if err != nil {
// 		ctx.Status(fiber.StatusNotFound)
// 		return
// 	}

// 	if body.Name != nil {
// 		todo.Name = *body.Name
// 	}

// 	if body.Completed != nil {
// 		todo.Completed = sql.NullBool{
// 			Bool:  *body.Completed,
// 			Valid: true,
// 		}
// 	}

// 	todo, err = h.Repo.UpdateTodo(ctx.Context(), postgres.UpdateTodoParams{
// 		ID:        int64(id),
// 		Name:      todo.Name,
// 		Completed: todo.Completed,
// 	})
// 	if err != nil {
// 		ctx.SendStatus(fiber.StatusNotFound)
// 		return
// 	}

// 	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
// 		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
// 		return
// 	}
// }

// func (h *Handlers) DeleteTodo(ctx *fiber.Ctx) {
// 	paramsId := ctx.Params("id")
// 	id, err := strconv.Atoi(paramsId)
// 	if err != nil {
// 		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "cannot parse id",
// 		})
// 		return
// 	}

// 	_, err = h.Repo.GetTodoById(ctx.Context(), int64(id))
// 	if err != nil {
// 		ctx.Status(fiber.StatusNotFound)
// 		return
// 	}

// 	err = h.Repo.DeleteTodoById(ctx.Context(), int64(id))
// 	if err != nil {
// 		ctx.SendStatus(fiber.StatusNotFound)
// 		return
// 	}

// 	ctx.SendStatus(fiber.StatusNoContent)
// }

func (h *Handlers) GetTodo(ctx *fiber.Ctx) {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	todo, err := h.Repo.GetTodoById(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}

func (h *Handlers) CreateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	if len(body.Name) <= 2 {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is not long enough",
		})
		return
	}

	todo, err := h.Repo.CreateTodo(ctx.Context(), postgres.CreateTodoParams{
		Name: body.Name,
	})

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(mapTodo(todo)); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}

func (h *Handlers) GetTodos(ctx *fiber.Ctx) {
	todos, err := h.Repo.GetAllTodos(ctx.Context())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	result := make([]interface{}, len(todos))
	for i, todo := range todos {
		result[i] = mapTodo(todo)
	}

	if err := ctx.Status(fiber.StatusOK).JSON(result); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}
