package middlewares

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net"
)

// OriginContextKey will be used as the request context key where the client's IP is stored
const (
	OriginContextKey string = "request-origin-addr"
)

func OriginFromRemoteAddr() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		c.Set(OriginContextKey, origin)

		c.Next()
	}
}

func OriginFromHeader(headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.Request.Header.Get(headerName)
		c.Set(OriginContextKey, headerValue)

		c.Next()
	}
}

func OriginFromBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer, _ := ioutil.ReadAll(io.LimitReader(c.Request.Body, 2048))
		origin := string(buffer)
		c.Set(OriginContextKey, origin)

		c.Next()
	}
}
