// Package middlewares provides common middleware handlers.
package middlewares

import (
	"errors"
	"net/http"

	"github.com/nstapelbroek/gatekeeper/libhttp"
	"github.com/spf13/viper"
)

// MustAuthenticate Enforces HTTP basic auth on a request and will respond early if the credentials do not match
func MustAuthenticate(config *viper.Viper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			username, password, success := libhttp.ParseBasicAuth(req.Header.Get("Authorization"))
			if !success {
				libhttp.BasicAuthUnauthorized(res, errors.New("Failed decoding basic auth header"))
				return
			}

			if config.GetString("http_auth_username") != username || config.GetString("http_auth_password") != password {
				libhttp.BasicAuthUnauthorized(res, errors.New("Username or password does not match"))
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
