package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound will add a generic json body to an HTTP 405 response
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
}

// MethodNotAllowed will add a generic json body to an HTTP 405 response
func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "HTTP verb not allowed"})
}
