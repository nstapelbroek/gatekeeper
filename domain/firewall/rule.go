package firewall

import (
	"net"
)

type Rule struct {
	Direction Direction
	Protocol  Protocol
	IP        net.IP
	Port      PortRange
}
