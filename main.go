package main

import (
	"log"

	"github.com/ADEMOLA200/hng-country-api/config"
	"github.com/ADEMOLA200/hng-country-api/repositories"
	"github.com/ADEMOLA200/hng-country-api/routes"
)

func main() {
	config.LoadConfig()

	err := repositories.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := routes.SetupRoutes()

	port := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(port))
}
