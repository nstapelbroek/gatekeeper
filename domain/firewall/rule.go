package firewall

import (
	"net"
)

// Rule is the value object used inside the application to create or delete firewall rules
type Rule struct {
	Direction Direction
	Protocol  Protocol
	IP        net.IP
	Port      PortRange
}
