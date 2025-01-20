package main

import (
	"archiv-system/internal/database"
	"archiv-system/internal/handler"
	"archiv-system/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Ensure the uploads directory exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatalf("Error creating uploads directory: %v", err)
	}

	// Initialize Gin
	r := gin.Default()

	// Initialize the database
	database.InitDB()

	// Public routes
	r.POST("auth/register", handler.Register)
	r.POST("auth/login", handler.Login)

	// Apply JWT middleware globally
	r.Use(middleware.JWTAuthMiddleware())

	// Group for document routes
	documentsGroup := r.Group("/documents")
	{
		documentsGroup.POST("/upload", middleware.AuthMiddleware("upload_document"), handler.UploadFile)
		documentsGroup.GET("/viewlist", middleware.AuthMiddleware("read_document"), handler.ViewListDoc) // Permission to view list of documents
		documentsGroup.PUT("/:id", middleware.AuthMiddleware("update_document"), middleware.OwnershipMiddleware(database.DB), handler.UpdateDocument)
		documentsGroup.DELETE("/:id", middleware.AuthMiddleware("delete_document"), middleware.OwnershipMiddleware(database.DB), handler.DeleteDocument)
		documentsGroup.GET("/user", middleware.AuthMiddleware("read_document"), handler.GetUserDocuments) // Permission to view user's own documents
		documentsGroup.GET("/:id/check-update", handler.CheckDocumentUpdate)
	}

	// Group for admin routes
	adminGroup := r.Group("/admin")
	{
		adminGroup.POST("/dashboard", middleware.AuthMiddleware("admin:create"), handler.AdminHandler)
	}

	// Group for user routes
	userGroup := r.Group("/user")
	{
		userGroup.GET("/dashboard", middleware.AuthMiddleware("user:read"), handler.UserHandler)
	}

	// Start the server
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
