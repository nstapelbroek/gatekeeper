package domain

import (
	"errors"
	"strconv"
	"strings"
)

// ErrStartHigherThanEnd is a constant used for saving the error that could occur when an invalid port number is selection
var (
	ErrStartHigherThanEnd = errors.New("the given start port number is higher than the end")
)

// PortRange is a value object containing the Port in a Firewall Rule
type PortRange struct {
	beginPort int
	endPort   int
}

// NewSinglePort is a constructor for a PortRange with only one port
func NewSinglePort(portnumber int) PortRange {
	var p PortRange
	p.beginPort = portnumber
	p.endPort = portnumber

	return p
}

// NewPortRange is a constructor for a PortRange with a port range
func NewPortRange(startPort int, endPort int) (PortRange, error) {
	var p PortRange

	if startPort > endPort {
		return p, ErrStartHigherThanEnd
	}

	p.beginPort = startPort
	p.endPort = endPort

	return p, nil
}

// NewPortFromString creates a portRange struct from an input string
func NewPortFromString(port string) (PortRange, error) {
	if !strings.Contains(port, "-") {
		port, err := strconv.ParseInt(port, 10, 0)

		return NewSinglePort(int(port)), err
	}

	portSlices := strings.Split(port, "-")
	startPort, _ := strconv.ParseInt(portSlices[0], 10, 0)
	endPort, _ := strconv.ParseInt(portSlices[1], 10, 0)

	return NewPortRange(int(startPort), int(endPort))
}

// IsSinglePort will evaluate if the PortRange contains a single port value
func (p PortRange) IsSinglePort() bool {
	return p.beginPort == p.endPort
}

// String will transform an PortRange to an string representation using a dash to separate begin and end port numbers
func (p PortRange) String() string {
	if p.IsSinglePort() {
		return strconv.Itoa(p.beginPort)
	}

	return strconv.Itoa(p.beginPort) + "-" + strconv.Itoa(p.endPort)
}
