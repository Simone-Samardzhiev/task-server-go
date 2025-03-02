package main

import (
	"log"
	"net/http"
	"server/config"
)

// App is the main entry of the application.
type App struct {
	config config.Config
}

// start will run the server.
func (a *App) start() error {
	server := http.Server{
		Addr: ":8080",
	}

	return server.ListenAndServe()
}

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		config: *conf,
	}

	err = app.start()
	if err != nil {
		log.Fatal(err)
	}
}
