package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(c *gin.Context) {
	var user models.User

	// Action to create user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if database.DB.Create(&user).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// action to hash password
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashpassword)

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var userLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		fmt.Println("Error in binding JSON:", err)
		return
	}

	fmt.Println("Received login request:", userLogin) // Log for debugging

	// Search user in the database
	var user models.User
	if err := database.DB.Where("username = ?", userLogin.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching user"})
		return
	}

	// Compare passwords
	if err := utils.ComparePassword(userLogin.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password incorrect"})
		return
	}

	// Generate the JWT token if everything is correct
	token, err := utils.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		fmt.Println("Token generation error:", err)
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"token": token})
}
