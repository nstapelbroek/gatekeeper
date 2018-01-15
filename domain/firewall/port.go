package firewall

import (
	"errors"
	"strconv"
)

var (
	ErrStartHigherThanEnd = errors.New("the given start port number is higher than the end")
)

type PortRange struct {
	beginPort int
	endPort   int
}

func NewSinglePort(portnumber int) (PortRange) {
	var p PortRange
	p.beginPort = portnumber
	p.endPort = portnumber

	return p
}

func NewPortRange(startPort int, endPort int) (PortRange, error) {
	var p PortRange

	if startPort > endPort {
		return p, ErrStartHigherThanEnd
	}

	p.beginPort = startPort
	p.endPort = endPort

	return p, nil
}

func (p PortRange) IsSinglePort() bool {
	return p.beginPort == p.endPort
}

func (p PortRange) String() string {
	if p.IsSinglePort() {
		return strconv.Itoa(p.beginPort)
	}

	return strconv.Itoa(p.beginPort) + "-" + strconv.Itoa(p.endPort)
}
