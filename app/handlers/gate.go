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
	adapterDispatcher *adapters.AdapterDispatcher
}

func NewGateHandler(timeoutConfig int64, rulesConfigValue string, dispatcher *adapters.AdapterDispatcher) (*gateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
	}

	h := gateHandler{
		defaultTimeout:    time.Duration(timeoutConfig) * time.Second,
		adapterDispatcher: dispatcher,
		defaultRules:      createRulesFromConfigString(rulesConfigValue),
	}

	return &h, nil
}

func (g gateHandler) PostOpen(c *gin.Context) {
	ipNet := g.getIpNetFromContext(c)
	rules := g.createRules(ipNet)

	_ = g.adapterDispatcher.Open(rules)
	//callResult := g.adapterDispatcher.Open(rules)
	//if !callResult.IsSuccessful() {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"error":   "Failed applying some rules",
	//		"details": "todo",
	//	})
	//	return
	//}

	timer := time.NewTimer(time.Duration(g.defaultTimeout))
	go func(rules []domain.Rule) {
		<-timer.C
		_ = g.adapterDispatcher.Close(rules)
	}(rules)


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
