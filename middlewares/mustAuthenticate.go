// Package middlewares provides common middleware handlers.
package middlewares

import (
	"errors"
	"net/http"

	"crypto/subtle"

	"github.com/nstapelbroek/gatekeeper/lib"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"
	"fmt"
)

// MustAuthenticate Enforces HTTP basic auth on a request and will respond early if the credentials do not match
func MustAuthenticate(config *viper.Viper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			username, password, success := lib.ParseBasicAuth(req.Header.Get("Authorization"))
			if !success {
				lib.BasicAuthUnauthorized(res, errors.New("Failed decoding basic auth header"))
				return
			}

			correctUsername := []byte(config.GetString("http_auth_username"))
			correctPassword := []byte(config.GetString("http_auth_password"))
			if subtle.ConstantTimeCompare(correctUsername, []byte(username)) == 0 || subtle.ConstantTimeCompare(correctPassword, []byte(password)) == 0 {
				logrus.Debugln(fmt.Sprintf("Authentication attempt failed, terminating request"))
				lib.BasicAuthUnauthorized(res, errors.New("Username or password does not match"))
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
