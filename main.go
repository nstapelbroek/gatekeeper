package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/adapters"
	"github.com/nstapelbroek/gatekeeper/handlers"
	"github.com/nstapelbroek/gatekeeper/middlewares"
	"github.com/spf13/viper"
)

func newConfig() *viper.Viper {
	c := viper.New()
	// Generic settings
	c.SetDefault("http_port", "8080")
	c.SetDefault("http_auth_username", "")
	c.SetDefault("http_auth_password", "")
	c.SetDefault("resolve_type", "RemoteAddr")
	c.SetDefault("resolve_header", "X-Forwarded-For")
	c.SetDefault("rule_close_timeout", 120)
	c.SetDefault("rule_ports", "TCP:22")

	// Adapter specific
	c.SetDefault("vultr_api_key", "")
	c.SetDefault("vultr_firewall_id", "")

	c.AutomaticEnv()

	return c
}

func registerBasicAuth(app *gin.Engine, config *viper.Viper) {
	username := config.GetString("http_auth_username")
	password := config.GetString("http_auth_password")
	if len(username) > 0 || len(password) > 0 {
		app.Use(gin.BasicAuth(gin.Accounts{username: password}))
	}
}

func registerResolver(app *gin.Engine, config *viper.Viper) {
	switch resolveType := config.GetString("resolve_type"); resolveType {
	case "headers":
		app.Use(middlewares.OriginFromHeader(config.GetString("resolve_header")))
	case "body":
		app.Use(middlewares.OriginFromBody())
	case "remote":
	default:
		app.Use(middlewares.OriginFromRemoteAddr())
	}
}

func registerRoutes(app *gin.Engine, config *viper.Viper) {
	adapterFactory := adapters.NewAdapterFactory(config)
	gateHandler := handlers.NewGateHandler(
		adapterFactory.GetAdapter(),
		config.GetInt("rule_close_timeout"),
		config.GetString("rule_ports"),
	)

	// Normal routes
	app.POST("/", gateHandler.PostOpen)

	// Error handling and helper routes
	app.Handle("GET", "/", handlers.MethodNotAllowed)
	app.Handle("PATCH", "/", handlers.MethodNotAllowed)
	app.Handle("PUT", "/", handlers.MethodNotAllowed)
	app.Handle("DELETE", "/", handlers.MethodNotAllowed)
	app.Handle("HEAD", "/", handlers.MethodNotAllowed)
	app.Handle("OPTIONS", "/", handlers.MethodNotAllowed)
	app.Handle("CONNECT", "/", handlers.MethodNotAllowed)
	app.Handle("TRACE", "/", handlers.MethodNotAllowed)
	app.NoRoute(handlers.NotFound)
}

func main() {
	config := newConfig()
	app := gin.Default()

	registerBasicAuth(app, config)
	registerResolver(app, config)
	registerRoutes(app, config)

	app.Run(":" + config.GetString("http_port"))
}
