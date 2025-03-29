package main

import (
	"auth/config"
	"auth/handlers"
	"auth/routes"
	"auth/storage"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	config.LoadConfig()
	cfg := config.GetConfig()

	// Initialize storage
	storageInstance, err := storage.NewStorage()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}
	defer storageInstance.Close()

	handlers.SetStorageInstance(storageInstance)

	r := routes.SetupRouter()
	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
