package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "HTTP verb not allowed"})
	return
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	return
}