package routes

import (
	"auth/auth"
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/signup", handlers.CreateProfile)
	r.POST("/signin", handlers.GetProfile)
	r.POST("/refresh", handlers.RefreshToken)

	// Protected routes
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.PUT("/forget-password", handlers.UpdateProfile)
		protected.POST("/delete-profile", handlers.DeleteProfile)
	}

	return r
}
