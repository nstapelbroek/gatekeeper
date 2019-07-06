package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
	"github.com/nstapelbroek/gatekeeper/domain"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
	"time"
)

type gateHandler struct {
	defaultTimeout    time.Duration
	defaultRules      []domain.Rule
	adapterDispatcher *adapters.AdapterDispatcher
	logger            *zap.Logger
}

func NewGateHandler(timeoutConfig int64, rulesConfigValue string, dispatcher *adapters.AdapterDispatcher, logger *zap.Logger) (*gateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
	}

	h := gateHandler{
		defaultTimeout:    time.Duration(timeoutConfig) * time.Second,
		adapterDispatcher: dispatcher,
		defaultRules:      createRulesFromConfigString(rulesConfigValue),
		logger:            logger,
	}

	return &h, nil
}

func (g gateHandler) PostOpen(c *gin.Context) {
	request := NewOpenRequestModel(c, g.defaultTimeout, g.defaultRules)
	g.logger.Debug(
		"Incoming open request",
		zap.String("ip", request.IpAddress.String()),
		zap.Duration("timeOut", request.Timeout),
	)

	r, err := g.adapterDispatcher.Open(request.Rules)
	if err != nil {
		g.openFailedResponse(c, err, r)
		return
	}

	g.openSuccessResponse(c, request.IpAddress, request.Timeout, r)

	if request.Timeout > 0 {
		g.scheduleDeletion(request)
	}
}

func (g gateHandler) openSuccessResponse(c *gin.Context, ip net.IPNet, timeout time.Duration, details map[string]string) {
	message := fmt.Sprintf("%s has been whitelisted", ip.String())

	if timeout > 0 {
		message = message + fmt.Sprintf("for %.0f seconds", timeout.Seconds())
	}

	c.JSON(http.StatusCreated, gin.H{"message": message, "details": details})
}

func (g gateHandler) openFailedResponse(c *gin.Context, err error, details map[string]string) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "details": details})
}

func (g gateHandler) scheduleDeletion(request *OpenRequestModel) {
	timer := time.NewTimer(time.Duration(g.defaultTimeout))
	go func(rules []domain.Rule) {
		<-timer.C
		g.logger.Debug("Closing", zap.String("ip", rules[0].IPNet.String()))
		_, _ = g.adapterDispatcher.Close(rules)
	}(request.Rules)
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
