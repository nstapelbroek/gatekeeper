package handlers

import (
	"net"
	"net/http"
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
	"github.com/nstapelbroek/gatekeeper/adapters/vultr"
)

func PostOpen(res http.ResponseWriter, req *http.Request) {
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

	adapter := vultr.Adapter{
		ApiKey:          "SomeKeyFromEnvironmentHere",
		FireWallGroupId: "SomeGroupIdFromEnvironmentHere",
	}

	adapter.CreateRule(rule)
	adapter.DeleteRule(rule)
}
