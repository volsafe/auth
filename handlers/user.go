package handlers

import (
	"auth/auth"
	"auth/storage"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

var S *storage.Storage

func SetStorageInstance(storageInstance *storage.Storage) {
	S = storageInstance
}

// User represents the user model for request/response
type User struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RefreshRequest represents the refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"your.refresh.token.here"`
}

// @Summary Create a new user profile
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User signup info"
// @Success 200 {object} map[string]string "message":"User created successfully"
// @Failure 400 {object} map[string]string "error":"Invalid JSON format"
// @Failure 500 {object} map[string]string "error":"Failed to create profile"
// @Router /signup [post]
func CreateProfile(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Error().Err(err).Msg("Failed to parse JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Info().Str("email", user.Email).Msg("Creating new user")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	nUser := storage.User{
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
	}

	err = S.SignUp(ctx, nUser)
	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	log.Info().Str("email", user.Email).Msg("User created successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// @Summary Sign in user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User credentials"
// @Success 200 {object} map[string]interface{} "tokens and user info"
// @Failure 400 {object} map[string]string "error":"Invalid JSON format"
// @Failure 401 {object} map[string]string "error":"Invalid credentials"
// @Failure 500 {object} map[string]string "error":"Failed to fetch profile"
// @Router /signin [post]
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

	// Check password
	err = S.CheckPassword(ctx, u.Email, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token pair
	tokenPair, err := auth.GenerateTokenPair(u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokenPair,
		"user":   profile,
	})
}

// @Summary Update user password
// @Description Update user's password (requires authentication)
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User info with new password"
// @Security BearerAuth
// @Success 200 {object} map[string]string "message":"Profile updated successfully"
// @Failure 400 {object} map[string]string "error":"Invalid JSON format"
// @Failure 401 {object} map[string]string "error":"Unauthorized"
// @Failure 500 {object} map[string]string "error":"Failed to update profile"
// @Router /forget-password [put]
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	nUser := storage.User{
		Email:        profile.Email,
		PasswordHash: string(hashedPassword),
	}

	err = S.UpdateUser(ctx, nUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	log.Info().Str("email", profile.Email).Msg("User profile updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// @Summary Delete user profile
// @Description Delete user's profile (requires authentication)
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User info"
// @Security BearerAuth
// @Success 200 {object} map[string]string "message":"Profile deleted successfully"
// @Failure 400 {object} map[string]string "error":"Invalid JSON format"
// @Failure 401 {object} map[string]string "error":"Unauthorized"
// @Failure 500 {object} map[string]string "error":"Failed to delete profile"
// @Router /delete-profile [post]
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

// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} auth.TokenPair "New token pair"
// @Failure 400 {object} map[string]string "error":"Invalid request format"
// @Failure 401 {object} map[string]string "error":"Invalid refresh token"
// @Router /refresh [post]
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate refresh token
	claims, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		if err == auth.ErrExpiredToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		}
		return
	}

	// Generate new token pair
	tokenPair, err := auth.GenerateTokenPair(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
