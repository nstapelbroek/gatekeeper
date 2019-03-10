package domain

import (
	"fmt"
	"net"
)

// Rule is the value object used inside the application to create or delete firewall rules
type Rule struct {
	Direction Direction
	Protocol  Protocol
	IPNet     net.IPNet
	Port      PortRange
}

func (r Rule) String() string {
	return fmt.Sprintf(
		"Rule-%s-%s-%s-%s-%s",
		r.Direction.String(),
		r.Protocol.String(),
		r.IPNet.IP.String(),
		r.IPNet.Mask.String(),
		r.Port.String(),
	)
}
