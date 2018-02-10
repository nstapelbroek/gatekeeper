package handlers

import (
	"net"
	"net/http"
	"time"
	"github.com/nstapelbroek/gatekeeper/adapters"
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
	"github.com/nstapelbroek/gatekeeper/middlewares"
)

type gateHandler struct {
	adapter adapters.Adapter
	timeout int
}

// NewGateHandler is an constructor for building gateHandler instances
func NewGateHandler(adapterInstance adapters.Adapter, timeout int) *gateHandler {
	h := new(gateHandler)
	h.adapter = adapterInstance
	h.timeout = timeout

	return h
}

func (handler gateHandler) PostOpen(res http.ResponseWriter, req *http.Request) {
	contextOrigin := req.Context().Value(middlewares.OriginContextKey)
	origin, assertionSucceeded := contextOrigin.(net.IP)
	if !assertionSucceeded {
		panic("context origin was somehow not the expected net.IP type")
	}

	direction, _ := firewall.NewDirectionFromString("inbound")
	protocol, _ := firewall.NewProtocolFromString("TCP")

	rule := firewall.Rule{
		Direction: direction,
		Protocol:  protocol,
		IP:        origin,
		Port:      firewall.NewSinglePort(22),
	}

	err := handler.adapter.CreateRule(rule)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Failed whitelisting, reason: " + err.Error()))
		return
	}

	timer := time.NewTimer(time.Second * 120)
	go func() {
		<-timer.C
		handler.adapter.DeleteRule(rule)
	}()

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(origin.String() + " has been whitelisted for 120 seconds"))
}
