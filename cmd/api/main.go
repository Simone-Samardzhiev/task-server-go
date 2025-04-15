package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"server/auth/tokens"
	"server/config"
	"server/database"
	"server/handlers"
	"server/repositories"
	"server/services"
)

type server struct {
	config        *config.Config
	handlers      handlers.Handlers
	authenticator *tokens.JWTAuthenticator
}

func (s *server) start() error {
	app := fiber.New()
	api := app.Group("/api")
	api1 := api.Group("/v1")

	// User routes
	userRouter := api1.Group("/users")
	userRouter.Post("/register", s.handlers.UserHandler.Register())
	userRouter.Post("/login", s.handlers.UserHandler.Login())
	userRouter.Get("/refresh", s.authenticator.Middleware(tokens.RefreshTokenType), s.handlers.UserHandler.Refresh())

	// Task routes
	taskRouter := api1.Group("/tasks", s.authenticator.Middleware(tokens.AccessTokenType))
	taskRouter.Get("/get", s.handlers.TaskHandler.GetTasks())
	taskRouter.Post("/add", s.handlers.TaskHandler.AddTask())
	taskRouter.Put("/update", s.handlers.TaskHandler.UpdateTask())
	taskRouter.Delete("/delete/:id", s.handlers.TaskHandler.DeleteTask())

	return app.Listen(s.config.ServerAddr)
}

func main() {
	conf := config.NewConfig()
	authenticator := tokens.NewJWTAuthenticator(&conf.AuthConfig)
	db, err := database.Connect(&conf.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error creating database connection: %v", err)
	}

	s := &server{
		authenticator: authenticator,
		config:        conf,
		handlers: handlers.Handlers{
			UserHandler: handlers.NewDefaultUserHandler(
				services.NewDefaultUserService(
					repositories.NewPostgresUserRepository(db),
					repositories.NewPostgresTokenRepository(db),
					authenticator,
				),
			),
			TaskHandler: handlers.NewDefaultTaskHandler(
				services.NewDefaultTaskService(
					repositories.NewPostgresTaskRepository(db),
				),
			),
		},
	}

	if err = s.start(); err != nil {
		log.Fatal("Error starting server", err)
	}
}
