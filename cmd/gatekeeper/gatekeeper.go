package main

import (
	"github.com/nstapelbroek/gatekeeper/app"
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
	c.SetDefault("app_env", "release")

	// Adapter specific
	c.SetDefault("vultr_api_key", "")
	c.SetDefault("vultr_firewall_group", "")

	c.AutomaticEnv()

	return c
}

func main() {
	a := app.NewApp(newConfig())

	_ = a.Run()
}
