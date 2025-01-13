package main

import (
	"log"
	"server/cmd/api"
	"server/config"
	"server/database"
)

func main() {
	db, err := database.CreateConnection(
		config.Envs.DBUser,
		config.Envs.DBPass,
		config.Envs.DBHost,
		config.Envs.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(config.Envs.Port, db)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
