package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

// testAccessLogMiddleware is a helper to prevent duplicated act logic
func testAccessLogMiddleware(t *testing.T, c *gin.Context, r *gin.Engine, hooks ...func(zapcore.Entry) error) {
	l := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(hooks...)))
	RegisterAccessLogMiddleware(r, l)

	r.HandleContext(c)
}

// setupGin is a helper to prevent duplicated arrange logic
func setupGin() (*gin.Context, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	return gin.CreateTestContext(httptest.NewRecorder())
}

func TestAccessLogRequestLineGet(t *testing.T) {
	c, r := setupGin()
	c.Request, _ = http.NewRequest(http.MethodGet, "/something-worth-asserting", nil)

	testAccessLogMiddleware(t, c, r, func(e zapcore.Entry) error {
		assert.Contains(t, e.Message, "\"GET /something-worth-asserting HTTP/1.1\"")
		return nil
	})
}

func TestAccessLogRequestLinePost(t *testing.T) {
	c, r := setupGin()
	c.Request, _ = http.NewRequest(http.MethodPost, "/something-worth-asserting", nil)

	testAccessLogMiddleware(t, c, r, func(e zapcore.Entry) error {
		assert.Contains(t, e.Message, "\"POST /something-worth-asserting HTTP/1.1\"")
		return nil
	})
}

func TestAccessLogAuthenticatedUser(t *testing.T) {
	// Since the order matters we'll do an additional test to make sure the assertion ran
	called := false

	c, r := setupGin()
	l := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
		assert.Contains(t, e.Message, " superuser [")
		called = true
		return nil
	})))

	RegisterAccessLogMiddleware(r, l)
	r.Use(gin.BasicAuth(gin.Accounts{"superuser": "password"}))
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	c.Request.Header.Add("Authorization", "Basic c3VwZXJ1c2VyOnBhc3N3b3Jk")

	r.HandleContext(c)
	assert.True(t, called)
}

func TestAccessLogTimeStamp(t *testing.T) {
	c, r := setupGin()
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)

	testAccessLogMiddleware(t, c, r, func(e zapcore.Entry) error {
		re := regexp.MustCompile(`(?m)\[[0-9]{2}\/[A-Z]{1}[a-z]{2}\/[0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} \+0000\]`)
		timeStamp := re.FindStringSubmatch(e.Message)
		assert.Len(t, timeStamp, 1)
		return nil
	})
}
