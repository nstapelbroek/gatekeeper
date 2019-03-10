package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
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
	adapters       []adapters.Adapter
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

func NewGateHandler(timeoutConfig int64, rulesConfigValue string, adapters []adapters.Adapter) (*gateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
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
	rules := make([]domain.Rule, len(g.defaultRules))
	copy(rules, g.defaultRules)
	var adapterErrors []error
	for _, rule := range rules {
		rule.IPNet = ipNet
		err := g.callAdapters(rule)
		if err != nil {
			adapterErrors = append(adapterErrors, err)
		}
	}

	if len(adapterErrors) > 0 {
		var details []string
		for _, adapterError := range adapterErrors {
			details = append(details, adapterError.Error())
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed applying some rules", "details": details})
		return
	}

	content := gin.H{
		"detail": fmt.Sprintf("%s has been whitelisted for %d seconds", ipNet.String(), g.defaultTimeout),
	}
	c.JSON(http.StatusCreated, content)
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

func (g gateHandler) callAdapters(rule domain.Rule) (err error) {
	// Todo loop over each adapter instead of picking the first one
	adapter := g.adapters[0]
	err = adapter.CreateRule(rule)

	if err == nil {
		timer := time.NewTimer(time.Duration(g.defaultTimeout))
		go func(rule domain.Rule) {
			<-timer.C
			_ = adapter.DeleteRule(rule)
		}(rule)
	}

	return
}
