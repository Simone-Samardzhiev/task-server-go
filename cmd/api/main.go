package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"server/config"
	"server/database"
)

func main() {
	conf := config.NewConfig()
	_, err := database.Connect(&conf.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	app := fiber.New()
	err = app.Listen(conf.ServerAddr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
