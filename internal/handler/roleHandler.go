package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AdminHandler(c *gin.Context) {
	// Initialisation des statistiques
	var totalUsers int64
	var totalAdmins int64
	var totalDocuments int64
	var activeUsers int64

	// Obtenir tous les utilisateurs
	if err := database.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch user statistics", nil)
		return
	}

	// Obtenir tous les admins
	if err := database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&totalAdmins).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch admin statistics", nil)
		return
	}

	// Obtenir tous les documents dans la base de données
	if err := database.DB.Model(&models.Document{}).Count(&totalDocuments).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch document statistics", nil)
		return
	}

	// Obtenir les utilisateurs actifs (ex. utilisateurs ayant téléchargé des documents au cours des 30 derniers jours)
	if err := database.DB.Model(&models.Document{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Distinct("user_id").
		Count(&activeUsers).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch active user statistics", nil)
		return
	}

	// Obtenir les documents par utilisateur
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
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch documents by user", nil)
		return
	}

	// Retourner les statistiques
	utils.RespondJSON(c, http.StatusOK, "Welcome to the Admin Dashboard", gin.H{
		"total_users":         totalUsers,
		"total_admins":        totalAdmins,
		"total_documents":     totalDocuments,
		"active_users":        activeUsers,
		"user_document_stats": userDocsCount,
	})
}

// UserHandler Logic
func UserHandler(c *gin.Context) {
	userID := c.GetString("userID")
	username := c.GetString("username")

	// Retourner les informations utilisateur
	utils.RespondJSON(c, http.StatusOK, "Welcome to the User Dashboard", gin.H{
		"id":       userID,
		"username": username,
	})
}
