package main

import (
	"github.com/nstapelbroek/gatekeeper/app"
	"github.com/spf13/viper"
)

func newConfig() *viper.Viper {
	c := viper.New()
	// Generic settings
	c.SetDefault("http_port", 8080)
	c.SetDefault("http_auth_username", "")
	c.SetDefault("http_auth_password", "")
	c.SetDefault("rule_close_timeout", 120)
	c.SetDefault("rule_ports", "TCP:22")
	c.SetDefault("app_env", "release")

	// Adapter specific
	c.SetDefault("digitalocean_personal_access_token", "")
	c.SetDefault("digitalocean_firewall_id", "")
	c.SetDefault("vultr_personal_access_token", "")
	c.SetDefault("vultr_firewall_id", "")
	c.SetDefault("aws_secret_key", "")
	c.SetDefault("aws_access_key", "")
	c.SetDefault("aws_region", "")
	c.SetDefault("aws_security_group_id", "")

	c.AutomaticEnv()

	return c
}

func main() {
	a := app.NewApp(newConfig())

	_ = a.Run()
}
