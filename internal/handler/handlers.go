package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashpassword)

	if database.DB.Create(&user).Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(201, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid input"})
		return
	}
	user := models.User{}
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logging Successfully", "token": token})
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload file"})
		return
	}

	destination := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	document := models.Document{
		Name: file.Filename,
		Type: file.Filename,
		URL:  destination,
		Tags: "", // I must add a logic to handle tags
	}
	if err := database.DB.Create(&document).Error; err != nil {
		return
	}
}

func ViewListDoc(c *gin.Context) {
	var documents []models.Document
	if database.DB.Find(&documents).Error != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch documents"})
		return
	}
	c.JSON(201, gin.H{"documents": documents})
}

// AdminHandler Logic
func AdminHandler(c *gin.Context) {
	// Tu peux inclure des statistiques, comme le nombre d'utilisateurs ou de documents.
	stats := gin.H{
		// fictitious data since I have not yet sent any requests to the database
		"total_users":     100,
		"total_documents": 500,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Admin Dashboard",
		"stats":   stats,
	})
}

// UserHandler Logic
func UserHandler(c *gin.Context) {
	userID := c.GetString("userID")
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the User Dashboard",
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
	})
}
