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
	defaultTimeout    int64
	defaultRules      []domain.Rule
	adapterDispatcher *adapters.AdapterDispatcher
	logger            *zap.Logger
}

type OpenRequestInput struct {
	Ip      string `form:"ip" json:"ip" binding:"omitempty,ip|cidr"`
	Timeout *int64 `form:"timeout" json:"timeout" binding:"omitempty,min=0,max=3600"`
}

func NewGateHandler(timeoutConfig int64, rulesConfigValue string, dispatcher *adapters.AdapterDispatcher, logger *zap.Logger) (*gateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
	}

	h := gateHandler{
		defaultTimeout:    timeoutConfig,
		adapterDispatcher: dispatcher,
		defaultRules:      createRulesFromConfigString(rulesConfigValue),
		logger:            logger,
	}

	return &h, nil
}

func (g gateHandler) PostOpen(c *gin.Context) {
	request, err := g.createOpenRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request", "details": err.Error()})
		return
	}

	g.logger.Debug(
		"Incoming open request",
		zap.String("ip", request.IpAddress.String()),
		zap.Duration("timeOut", request.Timeout),
	)
	r, err := g.adapterDispatcher.Open(request.Rules)
	if request.Timeout > 0 {
		g.scheduleDeletion(request)
	}

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error(), "details": r})
		return
	}

	g.openSuccessResponse(c, request.IpAddress, request.Timeout, r)
}

func (g gateHandler) scheduleDeletion(request *domain.OpenRequest) {
	timer := time.NewTimer(time.Duration(request.Timeout))
	go func(rules []domain.Rule) {
		<-timer.C
		g.logger.Debug("Closing", zap.String("ip", rules[0].IPNet.String()))
		_, _ = g.adapterDispatcher.Close(rules)
	}(request.Rules)
}

func (g gateHandler) openSuccessResponse(c *gin.Context, ip net.IPNet, timeout time.Duration, details map[string]string) {
	message := fmt.Sprintf("%s has been whitelisted", ip.String())

	if timeout > 0 {
		message = message + fmt.Sprintf("for %.0f seconds", timeout.Seconds())
	}

	c.JSON(http.StatusCreated, gin.H{"message": message, "details": details})
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

func (g gateHandler) createOpenRequest(c *gin.Context) (*domain.OpenRequest, error) {
	// Use model binding and validation from HTTP body. Fall back to origin IP when none is received
	var input OpenRequestInput
	if err := c.ShouldBind(&input); err != nil {
		return nil, err
	}

	if input.Timeout == nil {
		input.Timeout = &g.defaultTimeout
	}

	if input.Ip == "" {
		input.Ip = c.ClientIP()
	}

	if input.Ip == "" {
		return nil, errors.New("could not determine ip")
	}

	return domain.NewOpenRequest(input.Ip, *input.Timeout, g.defaultRules), nil
}
