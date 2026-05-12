package main

import (
	"crickxi-backend/database"
	"crickxi-backend/routes"
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

	srv := routes.SetUpRoutes()
	if err := srv.Run(":8080"); err != nil {
		panic(err)
	}
}
