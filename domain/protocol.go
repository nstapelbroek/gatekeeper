package domain

import (
	"errors"
	"strings"
)

// ErrInvalidProtocolString is a constant used for saving the error that could occur when an invalid protocol is passed
var (
	ErrInvalidProtocolString = errors.New("not a valid string to convert to a protocol")
)

// Protocol is an integer representation of the networking protocols used in a firewall rule
type Protocol int

// TCP is a constant value used in the Protocol value object
const (
	TCP Protocol = iota + 1
	UDP
	ICMP
)

// NewProtocolFromString is a constructor for Protocol
func NewProtocolFromString(protocol string) (Protocol, error) {
	var p Protocol

	switch strings.ToLower(protocol) {
	case "tcp", "6":
		p = TCP
	case "udp", "17":
		p = UDP
	case "icmp", "1":
		p = ICMP
	default:
		return p, ErrInvalidProtocolString
	}

	return p, nil
}

// String will convert the integer value of Protocol back to a capitalized string value
func (p Protocol) String() string {
	switch p {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case ICMP:
		return "icmp"
	default:
		return ""
	}
}

// ProtocolNumber will convert the object to an IANA protocol number, see https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
func (p Protocol) ProtocolNumber() int {
	switch p {
	case TCP:
		return 6
	case UDP:
		return 17
	case ICMP:
		return 1
	default:
		return -1
	}
}
