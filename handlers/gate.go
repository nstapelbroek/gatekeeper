package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/adapters"
	"github.com/nstapelbroek/gatekeeper/lib"
	"github.com/nstapelbroek/gatekeeper/middlewares"
	"net"
	"net/http"
	"time"
)

type gateHandler struct {
	adapter    adapters.Adapter
	timeout    int
	portConfig string
}

// NewGateHandler is an constructor for building gateHandler instances
func NewGateHandler(adapterInstance adapters.Adapter, timeout int, portConfig string) *gateHandler {
	h := new(gateHandler)
	h.adapter = adapterInstance
	h.timeout = timeout
	h.portConfig = portConfig

	return h
}

func (handler gateHandler) PostOpen(c *gin.Context) {
	contextOrigin := c.GetString(middlewares.OriginContextKey)
	originIP := net.ParseIP(contextOrigin)

	if originIP.IsLoopback() || originIP.IsLinkLocalUnicast() || originIP.IsLinkLocalMulticast() || originIP.IsLinkLocalUnicast() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Local IP's cannot be whitelisted"})
		return
	}

	timeOutInSeconds := int(time.Second) * handler.timeout
	rules := lib.CreateRules(handler.portConfig, originIP)

	for _, rule := range rules {
		err := handler.adapter.CreateRule(rule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		timer := time.NewTimer(time.Duration(timeOutInSeconds))
		go func() {
			<-timer.C
			handler.adapter.DeleteRule(rule)
		}()
	}

	content := gin.H{"detail": fmt.Sprintf("%s has been whitelisted for %d seconds", originIP.String(), handler.timeout)}
	c.JSON(http.StatusCreated, content)
}
