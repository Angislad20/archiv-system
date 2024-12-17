package main

import (
	"archiv-system/internal/database"
	"archiv-system/internal/handler"
	"archiv-system/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	database.InitDB()

	r.POST("auth/register", handler.Register)
	r.POST("auth/login", handler.Login)
	r.POST("/documents", handler.UploadFile)
	r.POST("/admin", middleware.AuthMiddleware("admin"), handler.AdminHandler)
	r.GET("user", middleware.AuthMiddleware("user"), handler.UserHandler)
	r.GET("/documents", handler.ViewListDoc)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
