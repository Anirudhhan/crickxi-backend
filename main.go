package main

import (
	"crickxi-backend/database"
)

func main() {

	err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"crickxi",
		"local",
		"local",
		database.SSLMode(database.SSLModeDisable),
	)

	if err != nil {
		panic(err)
	}
}
