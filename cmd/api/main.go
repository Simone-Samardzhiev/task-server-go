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
	mux.Handle("POST /users/register", a.handlers.UserHandler.Register())
	mux.Handle("POST /users/login", a.handlers.UserHandler.Login())
	mux.Handle(
		"GET /users/refresh",
		a.authenticator.Middleware(
			a.handlers.UserHandler.Refresh(),
			tokens.RefreshTokenType,
		),
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return server.ListenAndServe()
}

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(&conf.DatabaseConfig)
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
		},
		authenticator: authenticator,
	}

	err = app.start()
	if err != nil {
		log.Fatal(err)
	}
}
