package handlers

import (
	"profile/controllers"
	"github.com/gin-gonic/gin"

)



func HealthCheck(c *gin.Context) {
	err := controllers.HealthCheck(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
		}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
