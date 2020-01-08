package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// RegisterAccessLogMiddleware will setup request logging in Extended Log Format (ELF).
func RegisterAccessLogMiddleware(app *gin.Engine, logger *zap.Logger) {
	app.Use(LogRequests(logger))
}

// LogRequests is a wrapper around the handlerChain, writing an Extended Log Format entry on handled HTTP request
func LogRequests(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		sugar := logger.Sugar()

		sugar.Infof(
			"%s - %s [%s] \"%s %s HTTP/%d.%d\" %d %d \"%s\" \"%s\"",
			c.ClientIP(),
			c.GetString(gin.AuthUserKey),
			startTime.Format("02/Jan/1993:15:04:05 -0700"),
			c.Request.Method,
			path,
			c.Request.ProtoMajor,
			c.Request.ProtoMinor,
			c.Writer.Status(),
			c.Writer.Size(),
			c.Request.Referer(),
			c.Request.UserAgent(),
		)
	}
}
