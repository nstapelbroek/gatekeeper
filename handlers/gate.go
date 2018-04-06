package handlers

import (
	"fmt"
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

func (handler gateHandler) PostOpen(res http.ResponseWriter, req *http.Request) {
	contextOrigin := req.Context().Value(middlewares.OriginContextKey)
	origin, assertionSucceeded := contextOrigin.(net.IP)
	if !assertionSucceeded {
		panic("context origin was somehow not the expected net.IP type")
	}

	timeOutInSeconds := int(time.Second) * handler.timeout
	rules := lib.CreateRules(handler.portConfig, origin)

	for _, rule := range rules {
		err := handler.adapter.CreateRule(rule)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Failed whitelisting, reason: " + err.Error()))
			return
		}

		timer := time.NewTimer(time.Duration(timeOutInSeconds))
		go func() {
			<-timer.C
			handler.adapter.DeleteRule(rule)
		}()
	}

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s has been whitelisted for %d seconds", origin.String(), handler.timeout)))
}
