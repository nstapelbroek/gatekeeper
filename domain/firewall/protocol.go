package firewall

import (
	"strings"
	"errors"
)

// ErrInvalidProtocolString is a constant used for saving the error that could occur when an invalid protocol is passed
var (
	ErrInvalidProtocolString = errors.New("not a valid string to convert to a protocol")
)

// Protocol is an integer representation of the networking protocols used in a firewall rule
type Protocol int

// TCP is a constant value used in the Protocol value object
const (
	TCP  Protocol = iota + 1
	UDP
	ICMP
)

// NewProtocolFromString is a constructor for Protocol
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

// String will convert the integer value of Protocol back to a capitalized string value
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
