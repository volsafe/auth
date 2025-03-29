package main

import (
	"auth/config"
	"auth/handlers"
	"auth/routes"
	"auth/storage"
	"fmt"

	_ "auth/docs" // Import generated docs

	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Auth Service API
// @version         1.0
// @description     A JWT-based authentication service with user management capabilities.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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

	// Add Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
