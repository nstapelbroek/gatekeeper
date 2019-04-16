package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/middlewares"
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"net/http"
	"strings"
	"time"
)

type gateHandler struct {
	defaultTimeout time.Duration
	defaultRules   []domain.Rule
	adapters       []domain.Adapter
}

func NewGateHandler(timeoutConfig int64, rulesConfigValue string, adapters []domain.Adapter) (*gateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
	}

	if len(adapters) == 0 {
		return nil, errors.New("no adapters configured")
	}

	h := gateHandler{
		defaultTimeout: time.Duration(timeoutConfig) * time.Second,
		adapters:       adapters,
		defaultRules:   createRulesFromConfigString(rulesConfigValue),
	}

	return &h, nil
}

func (g gateHandler) PostOpen(c *gin.Context) {
	ipNet := g.getIpNetFromContext(c)
	rules := g.createRules(ipNet)

	errorDetails := make(map[string]string)
	for _, adapter := range g.adapters {
		callResult := adapter.CreateRules(rules)
		if !callResult.IsSuccessful() {
			errorDetails[adapter.ToString()] = callResult.Error.Error()
			continue
		}

		timer := time.NewTimer(time.Duration(g.defaultTimeout))
		go func(adapter domain.Adapter, rules []domain.Rule) {
			<-timer.C
			_ = adapter.DeleteRules(rules)
		}(adapter, rules)
	}

	if len(errorDetails) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Failed applying some rules",
			"details": errorDetails,
		})
		return
	}

	content := gin.H{
		"detail": fmt.Sprintf("%s has been whitelisted for %.0f seconds", ipNet.String(), g.defaultTimeout.Seconds()),
	}
	c.JSON(http.StatusCreated, content)
}

func (g gateHandler) createRules(ipNet net.IPNet) []domain.Rule {
	rules := make([]domain.Rule, len(g.defaultRules))
	copy(rules, g.defaultRules)
	for index := range rules {
		rules[index].IPNet = ipNet
	}

	return rules
}

func (g gateHandler) getIpNetFromContext(c *gin.Context) net.IPNet {
	contextOrigin := c.GetString(middlewares.OriginContextKey)
	if strings.Contains(contextOrigin, "/") {
		ip, ipNet, _ := net.ParseCIDR(contextOrigin)
		return net.IPNet{IP: ip, Mask: ipNet.Mask}

	}

	ip := net.ParseIP(contextOrigin)
	if ip.To4() != nil {
		return net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
	}

	// Assuming IPv6 here because conversion to IPv4 failed
	return net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
}

func createRulesFromConfigString(portConfig string) []domain.Rule {
	portsAndProtocols := strings.Split(portConfig, ",")
	rules := make([]domain.Rule, len(portsAndProtocols))

	for index, portAndProtocol := range portsAndProtocols {
		portAndProtocol := strings.SplitN(portAndProtocol, ":", 2)

		direction, _ := domain.NewDirectionFromString("inbound")
		protocol, _ := domain.NewProtocolFromString(portAndProtocol[0])
		port, _ := domain.NewPortFromString(portAndProtocol[1])

		newRule := domain.Rule{
			Direction: direction,
			Protocol:  protocol,
			Port:      port,
		}

		rules[index] = newRule
	}

	return rules
}
