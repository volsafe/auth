package main

import (
    "auth/handlers"
    "auth/routes"
    "auth/storage"
	_ "github.com/lib/pq"
	)

func main() {
    storageInstance, err := storage.NewStorage()
    if err != nil {
        panic("Failed to connect to the database")
    }
    defer storageInstance.Close()

    handlers.SetStorageInstance(storageInstance)

    r := routes.SetupRouter()
    r.Run(":8080")
}