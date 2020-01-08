package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// RegisterBasicAuthentication will setup gin basic auth middleware if a username and password is given
func RegisterBasicAuthentication(app *gin.Engine, config *viper.Viper) {
	username := config.GetString("http_auth_username")
	password := config.GetString("http_auth_password")
	if len(username) > 0 || len(password) > 0 {
		app.Use(gin.BasicAuth(gin.Accounts{username: password}))
	}
}
