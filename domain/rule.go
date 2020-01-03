package domain

import (
	"net"
	"strings"
)

// Rule is the value object used inside the application to create or delete firewall rules
type Rule struct {
	Direction Direction
	Protocol  Protocol
	IPNet     net.IPNet
	Port      PortRange
}

func (r Rule) String() string {
	return strings.Join(
		[]string{
			"Rule",
			r.Direction.String(),
			r.Protocol.String(),
			r.IPNet.IP.String(),
			r.IPNet.Mask.String(),
			r.Port.String(),
		}, "-")
}
