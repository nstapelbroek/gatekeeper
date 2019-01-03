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

	assert.Equal(t, "8080", config.GetString("http_port"))
	assert.Equal(t, "RemoteAddr", config.GetString("resolve_type"))
	assert.Equal(t, "X-Forwarded-For", config.GetString("resolve_header"))
	assert.Equal(t, 120, config.GetInt("rule_close_timeout"))
	assert.Equal(t, "TCP:22", config.GetString("rule_ports"))
	assert.Equal(t, "release", config.GetString("app_env"))
}
