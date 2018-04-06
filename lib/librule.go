package lib

import (
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
	"net"
	"strings"
)

func CreateRules(portConfig string, origin net.IP) []firewall.Rule {
	portsAndProtocols := strings.Split(portConfig, ",")
	rules := make([]firewall.Rule, len(portsAndProtocols))

	if len(portsAndProtocols) < 0 {
		return rules
	}

	for index, portAndProtocol := range portsAndProtocols {
		portAndProtocol := strings.SplitN(portAndProtocol, ":", 2)

		direction, _ := firewall.NewDirectionFromString("inbound")
		protocol, _ := firewall.NewProtocolFromString(portAndProtocol[0])
		port, _ := firewall.NewPortFromString(portAndProtocol[1])

		newRule := firewall.Rule{
			Direction: direction,
			Protocol:  protocol,
			IP:        origin,
			Port:      port,
		}

		rules[index] = newRule
	}

	return rules
}
