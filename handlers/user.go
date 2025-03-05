package handlers

import (
	"auth/storage"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

var S *storage.Storage

func SetStorageInstance(storageInstance *storage.Storage) {
	S = storageInstance
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}


func CreateProfile(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	nUser := storage.User{
		Email:        user.Email,
		PasswordHash: user.Password,
	}

	err = S.SignUp(ctx, nUser)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetProfile(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	profile, err := S.GetUser(ctx, u.Email)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()

	var profile User
	if err := json.Unmarshal(body, &profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	nUser := storage.User{
		Email:        profile.Email,
		PasswordHash: profile.Password,
	}

	err = S.UpdateUser(ctx, nUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func DeleteProfile(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err = S.DeleteUser(ctx, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
