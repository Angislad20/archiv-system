package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"time"
)

// Global variable with restricted access
var (
	jwtKey []byte
	once   sync.Once
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleName string `json:"roleName"`
	jwt.RegisteredClaims
}

// InitJWTKey initializes the jwtKey only once
func InitJWTKey() error {
	var initErr error
	once.Do(func() {
		// Load JWT_SECRET from environment variables
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			initErr = errors.New("JWT_SECRET is not set in environment variables")
			return
		}
		jwtKey = []byte(secret)
	})
	return initErr
}

// GenerateToken generates a JWT token for the user with the specified ID and role.
func GenerateToken(userID uint, roleName string) (string, error) {
	// Ensure jwtKey is initialized
	if err := InitJWTKey(); err != nil {
		return "", err
	}

	// Set expiration
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken verifies and extracts claims from the token
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure jwtKey is initialized
		if err := InitJWTKey(); err != nil {
			return nil, err
		}

		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

// ComparePassword securely compares hashed password and plain password
func ComparePassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
