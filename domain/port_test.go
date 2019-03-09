package domain

import "testing"

func TestNewSinglePort(t *testing.T) {
	port := NewSinglePort(123)

	if port.beginPort != 123 {
		t.Errorf("Port constructor failed to set the correct value")
	}

	if port.endPort != 123 {
		t.Errorf("Port constructor failed to set the correct value")
	}

	if !port.IsSinglePort() {
		t.Errorf("Port failed to detect that it was single")
	}

	if port.String() != "123" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "123")
	}
}

func TestNewPortRange(t *testing.T) {
	port, _ := NewPortRange(123, 125)

	if port.beginPort != 123 {
		t.Errorf("Port constructor failed to set the correct beginPort value")
	}

	if port.endPort != 125 {
		t.Errorf("Port constructor failed to set the correct endPort value")
	}

	if port.IsSinglePort() {
		t.Errorf("Port failed to detect that it was a range")
	}

	if port.String() != "123-125" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "123-125")
	}
}

func TestNewPortRangeInvalidValue(t *testing.T) {
	_, err := NewPortRange(20, 10)
	if err != ErrStartHigherThanEnd {
		t.Errorf("Portrange constructor did not return the expected error")
	}
}
