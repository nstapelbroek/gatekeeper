package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestDefaultConfigValues(t *testing.T) {
	config := newConfig()

	assert.Equal(t, 8080, config.GetInt("http_port"))
	assert.Equal(t, 120, config.GetInt("rule_close_timeout"))
	assert.Equal(t, "TCP:22", config.GetString("rule_ports"))
	assert.Equal(t, "release", config.GetString("app_env"))

	assert.Equal(t, "", config.GetString("digitalocean_personal_access_token"))
	assert.Equal(t, "", config.GetString("digitalocean_firewall_id"))
	assert.Equal(t, "", config.GetString("vultr_personal_access_token"))
	assert.Equal(t, "", config.GetString("vultr_firewall_id"))
	assert.Equal(t, "", config.GetString("aws_secret_key"))
	assert.Equal(t, "", config.GetString("aws_access_key"))
	assert.Equal(t, "", config.GetString("aws_region"))
	assert.Equal(t, "", config.GetString("aws_security_group_id"))
}
