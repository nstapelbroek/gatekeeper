package domain

import (
	"net"
)

// Rule is the value object used inside the application to create or delete firewall rules
type Rule struct {
	Direction Direction
	Protocol  Protocol
	IPNet     net.IPNet
	Port      PortRange
}
