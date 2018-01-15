package firewall

import (
	"strings"
	"errors"
)

var (
	ErrInvalidDirectionString = errors.New("not a valid string to convert to a direction")
)

type Direction int

const (
	Inbound  Direction = iota + 1
	Outbound
)

func NewDirectionFromString(direction string) (Direction, error) {
	var d Direction

	if strings.ToLower(direction) == "inbound" {
		d = Inbound
		return d, nil
	}

	if strings.ToLower(direction) == "outbound" {
		d = Outbound
		return d, nil
	}

	return d, ErrInvalidDirectionString
}

func (d Direction) IsInbound() bool {
	return d == Inbound
}

func (d Direction) IsOutbound() bool {
	return d == Outbound
}

func (d Direction) String() string {
	if d == Inbound {
		return "inbound"
	}

	if d == Outbound {
		return "outbound"
	}

	return ""
}
