package handlers

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
	"github.com/nstapelbroek/gatekeeper/domain"
	"go.uber.org/zap"
)

// GateHandler will have methods to handle HTTP requests for opening or closing rules
type GateHandler struct {
	defaultTimeout    int64
	defaultRules      []domain.Rule
	adapterDispatcher *adapters.AdapterDispatcher
	logger            *zap.Logger
}

type openRequestInput struct {
	IP      string `form:"ip" json:"ip" binding:"omitempty,ip|cidr"`
	Timeout *int64 `form:"timeout" json:"timeout" binding:"omitempty,min=0,max=3600"`
}

// NewGateHandler is a constructor for GateHandler
func NewGateHandler(timeoutConfig int64, rulesConfigValue string, dispatcher *adapters.AdapterDispatcher, logger *zap.Logger) (*GateHandler, error) {
	if len(rulesConfigValue) == 0 {
		return nil, errors.New("no rules configured")
	}

	h := GateHandler{
		defaultTimeout:    timeoutConfig,
		adapterDispatcher: dispatcher,
		defaultRules:      createRulesFromConfigString(rulesConfigValue),
		logger:            logger,
	}

	return &h, nil
}

// PostOpen will handle a HTTP Post request to open firewall rules
func (g GateHandler) PostOpen(c *gin.Context) {
	request, err := g.createOpenCommandFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request", "details": err.Error()})
		return
	}

	g.logger.Debug(
		"Incoming open request",
		zap.String("ip", request.IPAddress.String()),
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

	g.openSuccessResponse(c, request.IPAddress, request.Timeout, r)
}

func (g GateHandler) scheduleDeletion(request *domain.OpenCommand) {
	timer := time.NewTimer(time.Duration(request.Timeout))
	go func(rules []domain.Rule) {
		<-timer.C
		g.logger.Debug("Closing", zap.String("ip", rules[0].IPNet.String()))
		_, _ = g.adapterDispatcher.Close(rules)
	}(request.Rules)
}

func (g GateHandler) openSuccessResponse(c *gin.Context, ip net.IPNet, timeout time.Duration, details map[string]string) {
	message := fmt.Sprintf("%s has been whitelisted", ip.String())

	if timeout > 0 {
		message = message + fmt.Sprintf(" for %.0f seconds", timeout.Seconds())
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

func (g GateHandler) createOpenCommandFromContext(c *gin.Context) (*domain.OpenCommand, error) {
	// Use model binding and validation from HTTP body. Fall back to origin IP when none is received
	var input openRequestInput
	if err := c.ShouldBind(&input); err != nil {
		return nil, err
	}

	if input.Timeout == nil {
		input.Timeout = &g.defaultTimeout
	}

	if input.IP == "" {
		input.IP = c.ClientIP()
	}

	if input.IP == "" {
		return nil, errors.New("could not determine ip")
	}

	return domain.NewOpenCommand(input.IP, *input.Timeout, g.defaultRules), nil
}
