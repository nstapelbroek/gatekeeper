package controllers

import (
	"net"
	"net/http"
	"time"
	"github.com/nstapelbroek/gatekeeper/adapters"
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
)

type gateController struct {
	adapterFactory *adapters.AdapterFactory
	timeout        int
}

func NewGateController(factory *adapters.AdapterFactory, timeout int) *gateController {
	h := new(gateController)
	h.adapterFactory = factory
	h.timeout = timeout

	return h
}

func (handler gateController) PostOpen(res http.ResponseWriter, req *http.Request) {
	contextOrigin := req.Context().Value("origin")
	origin, assertionSuceeded := contextOrigin.(net.IP)
	if !assertionSuceeded {
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

	adapter := handler.adapterFactory.GetAdapter()
	err := adapter.CreateRule(rule)
	if err != nil {
		res.Write([]byte(err.Error()))
	}

	timer := time.NewTimer(time.Second * 10)
	go func() {
		<-timer.C
		adapter.DeleteRule(rule)
	}()
}