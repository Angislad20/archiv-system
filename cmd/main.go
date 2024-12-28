package main

import (
	"archiv-system/internal/database"
	"archiv-system/internal/handler"
	"archiv-system/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement des variables d'environnement : %v", err)
	}

	r := gin.Default()
	database.InitDB()

	r.POST("auth/register", handler.Register)
	r.POST("auth/login", handler.Login)
	r.POST("/documents", middleware.AuthMiddleware("create_document"), handler.UploadFile)
	r.GET("/documents", middleware.AuthMiddleware("read_document"), handler.ViewListDoc)
	r.POST("/admin", middleware.AuthMiddleware("admin"), handler.AdminHandler)
	r.GET("user", middleware.AuthMiddleware("user"), handler.UserHandler)
	r.PUT("/documents/:id", middleware.AuthMiddleware("update_document"), middleware.OwnershipMiddleware(database.DB), handler.UpdateDocument)
	r.DELETE("/documents/:id", middleware.AuthMiddleware("delete_document"), middleware.OwnershipMiddleware(database.DB), handler.DeleteDocument)

	r.GET("/documents/user", middleware.AuthMiddleware("read_document"), handler.GetUserDocuments)
	r.GET("/documents/:id/check-update", handler.CheckDocumentUpdate)

	err = r.Run(":8080")
	if err != nil {
		return
	}

}
