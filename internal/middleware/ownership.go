package middleware

import (
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func OwnershipMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get document ID from query parameters
		docIDStr := c.Param("id")
		docID, err := strconv.Atoi(docIDStr)
		if err != nil || docID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			c.Abort()
			return
		}

		var ownerID uint

		// Get logged-in user ID from context
		userID, exists := c.Get("UserID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		err = db.Table("documents").Select("owner_id").Where("id = ? AND owner_id = ?", docID, userID).Scan(&ownerID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this document or it does not exist"})
			c.Abort()
			return
		}
		// Check if user is the owner
		var doc models.Document
		if err := db.First(&doc, docID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			c.Abort()
			return
		}

		isOwner, err := utils.IsDocumentOwner(docID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check document ownership"})
			c.Abort()
			return
		}
		if !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this document"})
			c.Abort()
			return
		}
		c.Next()
	}
}
