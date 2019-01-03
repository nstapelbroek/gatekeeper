package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupGin(middleWare gin.HandlerFunc) *gin.Engine {
	app := gin.New()
	app.Use(middleWare)

	app.POST("/", func(c *gin.Context) {
		contextOrigin := c.GetString(OriginContextKey)
		c.String(200, contextOrigin)
	})

	return app
}

func performRequest(app http.Handler, request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
}

func TestOriginFromRemoteAddr(t *testing.T) {
	assertionIp := "123.123.123.123"
	middleware := OriginFromRemoteAddr()
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", nil)
	request.RemoteAddr = assertionIp + ":4567"

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}

func TestOriginFromBody(t *testing.T) {
	assertionIp := "127.0.0.1"
	middleware := OriginFromBody()
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(assertionIp)))

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}

func TestOriginFromHeader(t *testing.T) {
	assertionIp := "123.123.123.123"
	middleware := OriginFromHeader("x-forwarded-for")
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", nil)
	request.Header.Add("x-forwarded-for", assertionIp)

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}

func TestOriginFromHeaderCaseInsensitive(t *testing.T) {
	assertionIp := "123.123.123.213"
	middleware := OriginFromHeader("x-forwarded-for")
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", nil)
	request.Header.Add("X-FORWARDED-FOR", assertionIp)

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}

func TestOriginFromHeaderCustomHeader(t *testing.T) {
	assertionIp := "13.37.13.37"
	middleware := OriginFromHeader("some-header")
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", nil)
	request.Header.Add("some-header", assertionIp)

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}

// When receiving multiple headers the http library will return this first entry
func TestOriginFromHeaderDuplicateHeader(t *testing.T) {
	assertionIp := "12.3.4.9"
	middleware := OriginFromHeader("x-forwarded-for")
	app := setupGin(middleware)
	request, _ := http.NewRequest("POST", "/", nil)
	request.Header.Add("x-forwarded-for", assertionIp)
	request.Header.Add("x-forwarded-for", "127.0.0.1")
	request.Header.Add("x-forwarded-for", "10.0.0.2")

	response := performRequest(app, request)

	assert.Equal(t, assertionIp, response.Body.String())
}
