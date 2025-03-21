package main

import (
	"log"
	"net/http"
	"server/auth/tokens"
	"server/config"
	"server/database"
	"server/handlers"
	"server/repositories"
	"server/services"
	"server/utils"
)

// App is the main entry of the application.
type App struct {
	config        config.Config
	handlers      handlers.Handlers
	authenticator *tokens.JWTAuthenticator
}

// start will run the server.
func (a *App) start() error {
	mux := http.NewServeMux()
	mux.Handle("POST /users/register", utils.Logger(a.handlers.UserHandler.Register()))
	mux.Handle("POST /users/login", utils.Logger(a.handlers.UserHandler.Login()))
	mux.Handle(
		"GET /users/refresh",
		utils.Logger(
			a.authenticator.Middleware(
				a.handlers.UserHandler.Refresh(),
				tokens.RefreshTokenType,
			),
		),
	)
	mux.Handle(
		"GET /tasks",
		utils.Logger(
			a.authenticator.Middleware(
				a.handlers.TaskHandler.GetTasks(),
				tokens.AccessTokenType),
		),
	)
	mux.Handle(
		"POST /tasks",
		utils.Logger(
			a.authenticator.Middleware(
				a.handlers.TaskHandler.AddTask(),
				tokens.AccessTokenType),
		),
	)
	mux.Handle(
		"PUT /tasks",
		utils.Logger(
			a.authenticator.Middleware(
				a.handlers.TaskHandler.UpdateTask(),
				tokens.AccessTokenType),
		),
	)
	mux.Handle(
		"DELETE /tasks/{id}",
		utils.Logger(
			a.authenticator.Middleware(
				a.handlers.TaskHandler.DeleteTask(),
				tokens.AccessTokenType)),
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return server.ListenAndServe()
}

func main() {
	conf := config.NewConfig()
	db, err := database.Connect(&conf.DatabaseConfig)
	go utils.StartDeleteTokenTask(db)
	if err != nil {
		log.Fatal(err)
	}

	authenticator := tokens.NewJWTAuthenticator(conf.AuthConfig.JwtSecret, conf.AuthConfig.JwtIssuer)
	app := &App{
		config: *conf,
		handlers: handlers.Handlers{
			UserHandler: handlers.NewDefaultUserHandler(
				services.NewDefaultService(
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
		authenticator: authenticator,
	}

	err = app.start()
	if err != nil {
		log.Fatal(err)
	}
}
