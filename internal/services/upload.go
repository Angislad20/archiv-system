package services

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadFileInput struct {
	UserID uint
	File   *models.UploadedFile
	Tags   *models.Tag
}

// ProcessFileUpload handles the business logic for uploading a file
func ProcessFileUpload(input UploadFileInput) (*models.Document, error) {
	// Ensure the upload directory exists
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// Generate a unique filename
	uniqueFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), input.File.Filename)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save the file on the server
	if err := input.File.Save(filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Process tags
	if input.Tags == nil || input.Tags.Name == "" {
		input.Tags = &models.Tag{Name: "untagged"} // Tag par défaut
	}

	tagList, err := findOrCreateTags(strings.Split(input.Tags.Name, ","))
	if err != nil {
		return nil, fmt.Errorf("failed to process tags: %w", err)
	}

	// Create the document
	document := models.Document{
		Name:    input.File.Filename,
		Type:    input.File.ContentType,
		URL:     filePath,
		OwnerID: input.UserID,
		Tags:    &tagList,
	}

	// Save the document to the database
	if err := database.DB.Create(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	return &document, nil
}

// findOrCreateTags searches for or creates tags based on their names
func findOrCreateTags(tagNames []string) ([]models.Tag, error) {
	var tagList []models.Tag

	for _, tagName := range tagNames {
		tag := models.Tag{}
		if err := database.DB.Where("name = ?", tagName).First(&tag).Error; err != nil {
			// If the tag does not exist, create it
			tag = models.Tag{Name: tagName}
			if err := database.DB.Create(&tag).Error; err != nil {
				return nil, fmt.Errorf("failed to create tag '%s': %w", tagName, err)
			}
		}
		tagList = append(tagList, tag)
	}

	return tagList, nil
}
