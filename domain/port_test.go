package domain

import "testing"

func TestNewSinglePort(t *testing.T) {
	port := NewSinglePort(123)

	if port.BeginPort != 123 {
		t.Errorf("Port constructor failed to set the correct value")
	}

	if port.EndPort != 123 {
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

	if port.BeginPort != 123 {
		t.Errorf("Port constructor failed to set the correct BeginPort value")
	}

	if port.EndPort != 125 {
		t.Errorf("Port constructor failed to set the correct EndPort value")
	}

	if port.IsSinglePort() {
		t.Errorf("Port failed to detect that it was a range")
	}

	if port.String() != "123-125" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "123-125")
	}
}

func TestNewPortFromString(t *testing.T) {
	port, _ := NewPortFromString("8080")

	if port.BeginPort != 8080 {
		t.Errorf("Port constructor failed to set the correct BeginPort value")
	}

	if port.EndPort != 8080 {
		t.Errorf("Port constructor failed to set the correct EndPort value")
	}

	if !port.IsSinglePort() {
		t.Errorf("Port misinterpreted a single port as range")
	}

	if port.String() != "8080" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "8080")
	}
}

func TestNewPortFromStringWithWhiteSpaces(t *testing.T) {
	port, _ := NewPortFromString("   8080    ")

	if port.BeginPort != 8080 {
		t.Errorf("Port constructor failed to set the correct BeginPort value")
	}

	if port.EndPort != 8080 {
		t.Errorf("Port constructor failed to set the correct EndPort value")
	}

	if !port.IsSinglePort() {
		t.Errorf("Port misinterpreted a single port as range")
	}

	if port.String() != "8080" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "8080")
	}
}

func TestNewPortRangeFromStringWithWhiteSpaces(t *testing.T) {
	port, _ := NewPortFromString(" 20 - 22 ")

	if port.BeginPort != 20 {
		t.Errorf("Port constructor failed to set the correct BeginPort value")
	}

	if port.EndPort != 22 {
		t.Errorf("Port constructor failed to set the correct EndPort value")
	}

	if port.IsSinglePort() {
		t.Errorf("Port failed to detect that it was a range")
	}

	if port.String() != "20-22" {
		t.Errorf("Port to string conversion gone wrong , got %v want %v", port.String(), "20-22")
	}
}

func TestNewPortFromStringWithInvalidValue(t *testing.T) {
	_, err := NewPortFromString("20:22")
	if err.Error() != "strconv.ParseInt: parsing \"20:22\": invalid syntax" {
		t.Errorf("NewPortFromString did not fail to parse")
	}
}

func TestNewPortRangeInvalidValue(t *testing.T) {
	_, err := NewPortRange(20, 10)
	if err != ErrStartHigherThanEnd {
		t.Errorf("Portrange constructor did not return the expected error")
	}
}
