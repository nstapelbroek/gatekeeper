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

type contextKey int

// OriginContextKey will be used as the key of the request context where the client's IP is stored
const (
	OriginContextKey contextKey = iota
)

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
				origin, _, _ = net.SplitHostPort(req.RemoteAddr)
				logrus.Debugln(fmt.Sprintf("Origin's IP is set to %s via request.RemoteAddr", origin))
			}

			originIP := net.ParseIP(origin)
			if originIP == nil {
				logrus.Warningln("Request failed due to invalid resolvement")
				http.Error(res, "Failed resolving your IP", http.StatusInternalServerError)
				return
			}

			if originIP.IsLoopback() || originIP.IsLinkLocalUnicast() || originIP.IsLinkLocalMulticast() || originIP.IsLinkLocalUnicast() {
				logrus.Warningln("Terminated local request")
				http.Error(res, "Local IP's cannot be whitelisted", http.StatusUnprocessableEntity)
				return
			}

			req = req.WithContext(context.WithValue(req.Context(), OriginContextKey, originIP))

			next.ServeHTTP(res, req)
		})
	}
}
