package handler

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"archiv-system/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// Register permet de créer un nouvel utilisateur
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Récupération des données d'entrée
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	// Vérifier si le rôle "user" existe
	var role models.Role
	if err := database.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch default role", err.Error())
		log.Printf("Error fetching role 'user': %v", err)
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Créer l'utilisateur
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		RoleID:   role.ID, // Associer l'utilisateur au rôle "user"
	}

	// Enregistrer l'utilisateur dans la base de données
	if err := database.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	// Répondre avec succès
	utils.RespondJSON(c, http.StatusCreated, "User created successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Login permet à un utilisateur de se connecter
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Récupération des données d'entrée
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	// Recherche de l'utilisateur dans la base de données
	var user models.User
	if err := database.DB.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// Génération d'un token JWT
	token, err := utils.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	// Répondre avec succès
	utils.RespondJSON(c, http.StatusOK, "Login successful", gin.H{"token": token})
}
