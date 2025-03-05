package routes

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Profile routes
	r.POST("/signup", handlers.CreateProfile)
	r.POST("/signin", handlers.GetProfile)
	r.PUT("/forget-password", handlers.UpdateProfile)
	r.POST("/delete-profile", handlers.DeleteProfile)
	

	return r
}