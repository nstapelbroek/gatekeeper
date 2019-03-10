package domain

import (
	"testing"
)

func testNewDirectionFromString(testValue string, t *testing.T) {
	direction, _ := NewDirectionFromString(testValue)
	if direction.String() != testValue {
		t.Errorf("Direction returned the wrong string value, got %v want %v", direction.String(), testValue)
	}
}

func TestNewDirectionFromStringInbound(t *testing.T) {
	testNewDirectionFromString("inbound", t)
}

func TestNewDirectionFromStringOutbound(t *testing.T) {
	testNewDirectionFromString("outbound", t)
}

func TestNewDirectionFromStringGracefullyHandlesCapitalisation(t *testing.T) {
	direction, _ := NewDirectionFromString("INBounD")
	if direction.String() != "inbound" {
		t.Errorf("Direction returned the wrong string value, got %v want %v", direction.String(), "inbound")
	}
}

func TestNewDirectionFromInvalidString(t *testing.T) {
	_, err := NewDirectionFromString("This_value_sucks")
	if err != ErrInvalidDirectionString {
		t.Errorf("Direction constructor did not return the expected error")
	}
}

func TestDirection_IsInbound(t *testing.T) {
	direction, _ := NewDirectionFromString("INBOUND")
	if !direction.IsInbound() {
		t.Errorf("Direction failed to identify itself as inbound")
	}
}

func TestDirection_IsOutbound(t *testing.T) {
	direction, _ := NewDirectionFromString("OUTBOUND")
	if !direction.IsOutbound() {
		t.Errorf("Direction failed to identify itself as outbound")
	}
}
