package handlers

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupGin(handler gin.HandlerFunc) *gin.Engine {
	app := gin.New()
	app.GET("/", handler)

	return app
}

func performRequest(app http.Handler, request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
}
