package main

import (
	"crickxi-backend/database"
	"log"
	"os"
)

func main() {

	err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"todo",
		"local",
		"local",
		database.SSLMode(database.SSLModeDisable),
	)

	if err != nil {
		panic(err)
	}

	secret, exists := os.LookupEnv("ACCESS_SECRET")
	if !exists || secret == "" {
		log.Fatal("ACCESS_SECRET is not set")
	}
}
