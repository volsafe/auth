package auth

import (
	"encoding/base64"
	"errors"
	"os"
	"time"
	"io"

	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenSecret  = loadSecret("ACCESS_TOKEN_SECRET")
	RefreshTokenSecret = loadSecret("REFRESH_TOKEN_SECRET")

	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// loadSecret loads and decodes a base64 encoded secret from environment variables
func loadSecret(key string) []byte {
	secret := os.Getenv(key)
	if secret == "" {
		panic("Missing required environment variable: " + key)
	}

	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		panic("Invalid base64 encoded secret for: " + key)
	}

	return decoded
}

// GenerateSecretString generates a base64 encoded secret string
func GenerateSecretString(length int) string {
	secret := make([]byte, length)
	file, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := io.ReadFull(file, secret); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(secret)
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(email string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := generateAccessToken(email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(email)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateAccessToken creates a new short-lived JWT token
func generateAccessToken(email string) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

// generateRefreshToken creates a new long-lived refresh token
func generateRefreshToken(email string) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token expires in 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshTokenSecret)
}

// ValidateToken validates the access token and returns the claims if valid
func ValidateToken(tokenString string, tokenType string) (*Claims, error) {
	switch tokenType {
	case "access":
		return validateToken(tokenString, AccessTokenSecret)
	case "refresh":
		return validateToken(tokenString, RefreshTokenSecret)
	default:
		var ErrInvalidTokenType = errors.New("invalid token type")
		return nil, ErrInvalidTokenType
	}
}

// ValidateRefreshToken validates the refresh token and returns the claims if valid
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, RefreshTokenSecret)
}

// validateToken is a helper function to validate tokens
func validateToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
