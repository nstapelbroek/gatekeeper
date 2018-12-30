package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMethodNotAllowed(t *testing.T) {
	app := setupGin(MethodNotAllowed)
	request, _ := http.NewRequest("GET", "/", nil)
	response := performRequest(app, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("content-type"))
	assert.Contains(t, string(response.Body.Bytes()), "HTTP verb not allowed")
}

func TestNotFound(t *testing.T) {
	app := setupGin(NotFound)
	request, _ := http.NewRequest("GET", "/", nil)
	response := performRequest(app, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("content-type"))
	assert.Contains(t, string(response.Body.Bytes()), "Page not found")
}