package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondJSON sends a JSON response with a status, message, and data
func RespondJSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"status":  http.StatusText(status),
		"message": message,
		"data":    data,
	})
}

// RespondError sends a JSON response with a status, message, and errors
func RespondError(c *gin.Context, status int, message string, errors interface{}) {
	c.JSON(status, gin.H{
		"status":  http.StatusText(status),
		"message": message,
		"errors":  errors,
	})
}
