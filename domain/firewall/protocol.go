package firewall

import (
	"strings"
	"errors"
)

var (
	ErrInvalidProtocolString = errors.New("not a valid string to convert to a direction")
)

type Protocol int

const (
	TCP  Protocol = iota + 1
	UDP
	ICMP
)

func NewProtocolFromString(protocol string) (Protocol, error) {
	var p Protocol

	if strings.ToLower(protocol) == "tcp" {
		p = TCP
		return p, nil
	}

	if strings.ToLower(protocol) == "udp" {
		p = UDP
		return p, nil
	}

	if strings.ToLower(protocol) == "icmp" {
		p = ICMP
		return p, nil
	}

	return p, ErrInvalidProtocolString
}

func (p Protocol) String() string {
	switch p {
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	case ICMP:
		return "ICMP"
	default:
		return ""
	}
}
