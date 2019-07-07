package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RegisterBasicAuthentication(app *gin.Engine, config *viper.Viper) {
	username := config.GetString("http_auth_username")
	password := config.GetString("http_auth_password")
	if len(username) > 0 || len(password) > 0 {
		app.Use(gin.BasicAuth(gin.Accounts{username: password}))
	}
}
