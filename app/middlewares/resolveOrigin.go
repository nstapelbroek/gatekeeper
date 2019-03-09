package middlewares

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net"
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
