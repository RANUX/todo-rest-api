package apiserver

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq" // ...
	"github.com/ranux/todo-rest-api/internal/app/postgres"
)

// Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// store := sqlstore.New(db)
	// sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	// srv := newServer(store, sessionStore)

	repo := postgres.NewRepo(db)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	handlers := NewHandlers(repo)

	SetupApiV1(app, handlers)

	return app.Listen(config.BindAddr)
}

func SetupApiV1(app *fiber.App, handlers *Handlers) {
	v1 := app.Group("/api/v1")

	SetupTodosRoutes(v1, handlers)
}

func SetupTodosRoutes(grp fiber.Router, handlers *Handlers) {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", handlers.GetTodos)
	todosRoutes.Post("/create/", handlers.CreateTodo)
	todosRoutes.Get("/:id", handlers.GetTodo)
	todosRoutes.Delete("/delete/:id", handlers.DeleteTodo)
	todosRoutes.Patch("/update/:id", handlers.UpdateTodo)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
