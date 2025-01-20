package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/services"
	"archiv-system/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile handles the uploading of documents
func UploadFile(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "Invalid User ID type", nil)
		return
	}

	// Récupérer le fichier de la requête
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Failed to get file", err.Error())
		return
	}

	// Construire un objet représentant le fichier
	uploadedFile := &models.UploadedFile{
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Save: func(destination string) error {
			return c.SaveUploadedFile(file, destination)
		},
	}

	// Récupérer les tags de la requête
	tags := c.PostForm("tags")

	// Appeler la logique métier
	input := services.UploadFileInput{
		UserID: userIDUint,
		File:   uploadedFile,
		Tags:   &models.Tag{Name: tags},
	}

	document, err := services.ProcessFileUpload(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Répondre avec succès
	utils.RespondJSON(c, http.StatusOK, "File uploaded successfully", gin.H{
		"ID":        document.ID,
		"Name":      document.Name,
		"Type":      document.Type,
		"URL":       document.URL,
		"Tags":      document.Tags,
		"CreatedAt": document.CreatedAt,
		"UpdatedAt": document.UpdatedAt,
	})
}

// ViewListDoc handles the retrieval of all documents
func ViewListDoc(c *gin.Context) {
	var documents []models.Document
	if database.DB.Find(&documents).Error != nil {
		utils.RespondError(c, http.StatusBadRequest, "Failed to fetch documents", gin.H{"error": "Failed to fetch documents"})
		return
	}
	utils.RespondJSON(c, http.StatusCreated, "Documents fetched successfully", gin.H{"documents": documents})
}

// UpdateDocument handles the updating of a document
func UpdateDocument(c *gin.Context) {
	docID := c.Param("id")
	//valider les entrées
	var updateRequest models.UpdateRequest
	if err := c.BindJSON(&updateRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Appeler la logique métier
	documentService := &services.DocumentService{}
	updateDocument, err := documentService.ProcessFileUpdate(docID, updateRequest)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update document", err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Document updated successfully", gin.H{"document": updateDocument})
}

// GetUserDocuments handles the retrieval of documents for a specific user
func GetUserDocuments(c *gin.Context) {
	userID := c.GetUint("UserID") // Retrieve the user ID from the context

	documents, err := utils.GetDocumentsByOwnerID(database.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch documents", err.Error())
		return
	}
	utils.RespondJSON(c, http.StatusOK, "Documents fetched successfully", gin.H{"documents": documents})
}

// GetDocumentsByTags handles the retrieval of documents by tags
func GetDocumentsByTags(c *gin.Context) {
	tags := c.DefaultQuery("tags", "") // Retrieve tags from the query (comma-separated)

	if tags == "" {
		utils.RespondError(c, http.StatusBadRequest, "Tags query parameter is required", nil)
		return
	}

	tagNames := strings.Split(tags, ",")
	var documents []models.Document

	// Search for documents associated with the given tags
	if err := database.DB.Preload("Tags").Joins("JOIN document_tags ON document_tags.document_id = documents.id").
		Joins("JOIN tags ON tags.id = document_tags.tag_id").
		Where("tags.name IN (?)", tagNames).
		Find(&documents).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to find documents by tags", err.Error())
		return
	}

	// Respond with the found documents
	utils.RespondJSON(c, http.StatusOK, "Documents retrieved successfully", gin.H{
		"documents": documents,
	})
}

// DeleteDocument handles the deletion of a document
func DeleteDocument(c *gin.Context) {
	docID := c.Param("id")

	var document models.Document
	if err := database.DB.First(&document, docID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Document not found", gin.H{"error": "Document not found"})
		return
	}

	if err := database.DB.Delete(&document).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete document", err.Error())
		return
	}
	utils.RespondJSON(c, http.StatusOK, "Document deleted successfully", gin.H{"document": document})
}

// CheckDocumentUpdate checks if a document has been updated since it was last viewed
func CheckDocumentUpdate(c *gin.Context) {
	docID := c.Param("id")
	lastViewedTimeStr := c.GetHeader("LastViewed")
	lastViewedTime, err := time.Parse(time.RFC3339, lastViewedTimeStr)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid timestamp", err.Error())
		return
	}

	updatedAt, err := utils.GetDocumentUpdatedAt(docID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to check update status", err.Error())
		return
	}

	if updatedAt.After(lastViewedTime) {
		utils.RespondJSON(c, http.StatusOK, "Document updated", gin.H{"update_available": true})
	} else {
		utils.RespondJSON(c, http.StatusOK, "Document not updated", gin.H{"update_available": false})
	}
}
