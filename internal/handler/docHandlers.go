package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// UploadFile Logic
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	// Validate file format
	if !strings.HasSuffix(file.Filename, ".pdf") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Define destination and save the file
	destination := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create document entry in the database
	document := models.Document{
		Name: file.Filename,
		Type: "pdf",
		URL:  destination,
		Tags: "", // Logic for tags can be added later
	}
	if err := database.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document in database"})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully", "document": document})
}

func ViewListDoc(c *gin.Context) {
	var documents []models.Document
	if database.DB.Find(&documents).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch documents"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"documents": documents})
}

func UpdateDocument(c *gin.Context) {
	var document models.Document
	docID := c.Param("id")
	if err := database.DB.First(&document, docID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if err := c.BindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := database.DB.Save(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func GetUserDocuments(c *gin.Context) {
	userID := c.GetUint("UserID") // Récupère l'ID de l'utilisateur connecté depuis le contexte

	documents, err := utils.GetDocumentsByOwnerID(database.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func DeleteDocument(c *gin.Context) {
	docID := c.Param("id")

	var document models.Document
	if err := database.DB.First(&document, docID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if err := database.DB.Delete(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

func CheckDocumentUpdate(c *gin.Context) {
	docID := c.Param("id")
	lastViewedTimeStr := c.GetHeader("LastViewed")
	lastViewedTime, err := time.Parse(time.RFC3339, lastViewedTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
		return
	}

	updatedAt, err := utils.GetDocumentUpdatedAt(docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check update status"})
		return
	}

	if updatedAt.After(lastViewedTime) {
		c.JSON(http.StatusOK, gin.H{"update_available": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"update_available": false})
	}
}
