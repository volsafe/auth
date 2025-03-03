package handlers

import (
	"auth/storage"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

var S *storage.Storage

func SetStorageInstance(storageInstance *storage.Storage) {
	S = storageInstance
}

type User struct {
	ID             int        `json:"id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"password_hash"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastLogin      time.Time  `json:"last_login,omitempty"`
	IsActive       bool       `json:"is_active"`
	FailedAttempts int        `json:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
}

func CreateProfile(c *gin.Context) {
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

	newProfile := storage.Profile{
		UserID:             profile.UserID,
		PhoneNumber:        profile.PhoneNumber,
		Latitude:           profile.Latitude,
		Longitude:          profile.Longitude,
		HomeSize:           profile.HomeSize,
		BuildingFloors:     profile.BuildingFloors,
		FloorLivingOn:      profile.FloorLivingOn,
		WindowModel:        profile.WindowModel,
		AdultsCount:        profile.AdultsCount,
		ChildrenCount:      profile.ChildrenCount,
		ElectricityCompany: profile.ElectricityCompany,
		MeterType:          profile.MeterType,
		BillNumber:         profile.BillNumber,
	}

	err = S.CreateUserProfile(ctx, newProfile)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func GetProfile(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := S.GetUserProfile(ctx, uint(userID))
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

	Profile := storage.Profile{
		UserID:             profile.UserID,
		PhoneNumber:        profile.PhoneNumber,
		Latitude:           profile.Latitude,
		Longitude:          profile.Longitude,
		HomeSize:           profile.HomeSize,
		BuildingFloors:     profile.BuildingFloors,
		FloorLivingOn:      profile.FloorLivingOn,
		WindowModel:        profile.WindowModel,
		AdultsCount:        profile.AdultsCount,
		ChildrenCount:      profile.ChildrenCount,
		ElectricityCompany: profile.ElectricityCompany,
		MeterType:          profile.MeterType,
		BillNumber:         profile.BillNumber,
	}

	err = S.UpdateUserProfileByID(ctx, Profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func DeleteProfile(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = S.DeleteUserProfile(ctx, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
