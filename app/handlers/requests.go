package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"strings"
	"time"
)

type OpenRequestModel struct {
	Rules           []domain.Rule
	Timeout         time.Duration
	IpAddress       net.IPNet
	RawIpValue      string `form:"ip" json:"ip" binding:"required"`
	RawTimeoutValue int    `form:"timeout" json:"timeout"`
}

func NewOpenRequestModel(c *gin.Context, defaultTimeout time.Duration, defaultRules []domain.Rule) *OpenRequestModel {
	// Use model binding and validation from HTTP body. Fall back to origin IP when none is received
	var or OpenRequestModel
	if err := c.ShouldBind(&or); err != nil {
		or.RawIpValue = c.ClientIP()
	}

	if or.RawIpValue == "" {
		// todo: handle better
		panic("oof! https://www.youtube.com/watch?v=aQU8LzYIxBU")
	}

	or.IpAddress = parseIp(or.RawIpValue)
	or.Timeout = parseTimeout(or.RawTimeoutValue, defaultTimeout)
	or.Rules = parseRules(or.IpAddress, defaultRules)

	return &or
}

func parseIp(stringIp string) net.IPNet {
	if strings.Contains(stringIp, "/") {
		ip, ipNet, _ := net.ParseCIDR(stringIp)
		return net.IPNet{IP: ip, Mask: ipNet.Mask}
	}

	ip := net.ParseIP(stringIp)
	if ip.To4() != nil {
		return net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
	}

	// Assuming IPv6 here because conversion to IPv4 failed
	return net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
}

func parseTimeout(requestedTimeout int, defaultTimeout time.Duration) time.Duration {
	// todo overwrite by request support
	return defaultTimeout
}

func parseRules(ipAddress net.IPNet, defaultRules []domain.Rule) []domain.Rule {
	rules := make([]domain.Rule, len(defaultRules))
	copy(rules, defaultRules)
	for index := range rules {
		rules[index].IPNet = ipAddress
	}

	return rules
}
