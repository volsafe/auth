package routes

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Profile routes
	r.POST("/signin", handlers.CreateProfile)
	r.POST("/signout", handlers.GetProfile)
	r.PUT("//forget-password", handlers.UpdateProfile)
	r.DELETE("/delete-profile", handlers.DeleteProfile)
	

	return r
}