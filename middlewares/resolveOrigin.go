// Package middlewares provides common middleware handlers.
package middlewares

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// OriginContextKey will be used as the key of the request context where the client's IP is stored
const OriginContextKey = "origin"

// ResolveOrigin Enforces HTTP basic auth on a request and will respond early if the credentials do not match
func ResolveOrigin(config *viper.Viper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			var origin string

			switch resolveType := config.GetString("resolve_type"); resolveType {
			case "headers":
				targetHeader := config.GetString("resolve_header")
				origin = req.Header.Get(targetHeader)
				logrus.Debugln(fmt.Sprintf("Origin's IP is set to %s via HTTP header %s", origin, targetHeader))
			case "body":
				body, err := ioutil.ReadAll(req.Body)
				origin = string(body)
				logrus.Debugln(fmt.Sprintf("Origin's IP is set to %s via Request body, errors: %s", origin, err))
			default:
				origin = req.RemoteAddr
				logrus.Debugln(fmt.Sprintf("Origin's IP is set to %s via request.RemoteAddr", origin))
			}

			if origin == "" {
				logrus.Warningln("Request failed due to invalid resolvement")
				http.Error(res, "Failed resolving your IP", http.StatusInternalServerError)
			}
			req = req.WithContext(context.WithValue(req.Context(), OriginContextKey, net.ParseIP(origin)))

			next.ServeHTTP(res, req)
		})
	}
}
