package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// AdminHandler Logic
func AdminHandler(c *gin.Context) {
	// initialize the statistics
	var totalUsers int64
	var totalAdmins int64
	var totalDocuments int64
	var activeUsers int64

	// Get all users
	if err := database.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user statistics"})
		return
	}

	// Get all admins
	if err := database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&totalAdmins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin statistics"})
		return
	}

	// fetch all docs in the database
	if err := database.DB.Model(&models.Document{}).Count(&totalDocuments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document statistics"})
		return
	}

	// fetch all actifs users (for example, users who download documents during the last 30 days)
	if err := database.DB.Model(&models.Document{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Distinct("user_id").
		Count(&activeUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active user statistics"})
		return
	}

	// fetch documents by user
	var userDocsCount []struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		DocCount int64  `json:"doc_count"`
	}
	if err := database.DB.Table("users").
		Select("users.id AS user_id, users.username, COUNT(documents.id) AS doc_count").
		Joins("LEFT JOIN documents ON users.id = documents.user_id").
		Group("users.id, users.username").
		Order("doc_count DESC").
		Limit(10).
		Scan(&userDocsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents by user"})
		return
	}

	// Return statistics
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Admin Dashboard",
		"stats": gin.H{
			"total_users":     totalUsers,
			"total_admins":    totalAdmins,
			"total_documents": totalDocuments,
			"active_users":    activeUsers,
		},
		"user_document_stats": userDocsCount,
	})
}

// UserHandler Logic
func UserHandler(c *gin.Context) {
	userID := c.GetString("userID")
	username := c.GetString("username")

	// Return statistics
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the User Dashboard",
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
	})
}
