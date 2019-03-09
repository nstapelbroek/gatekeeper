package domain

import (
	"errors"
	"strings"
)

// ErrInvalidDirectionString is a constant used for saving the error that could occur when an invalid direction is configured
var (
	ErrInvalidDirectionString = errors.New("not a valid string to convert to a direction")
)

// Direction is an integer representation (value object) of "inbound" or "outbound"
type Direction int

// Inbound is a constant value used in the Direction value object
const (
	Inbound Direction = iota + 1
	Outbound
)

// NewDirectionFromString constructs a Direction from string
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

// IsInbound evaluates of the Direction is "inbound"
func (d Direction) IsInbound() bool {
	return d == Inbound
}

// IsOutbound evaluates of the Direction is "outbound"
func (d Direction) IsOutbound() bool {
	return d == Outbound
}

// String will convert the integer back to an string
func (d Direction) String() string {
	if d == Inbound {
		return "inbound"
	}

	if d == Outbound {
		return "outbound"
	}

	return ""
}
