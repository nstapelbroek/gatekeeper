package domain

import (
	"testing"
)

func testNewProtocolFromString(testValue string, t *testing.T) {
	protocol, _ := NewProtocolFromString(testValue)
	if protocol.String() != testValue {
		t.Errorf("Protocol returned the wrong string value, got %v want %v", protocol.String(), testValue)
	}
}

func TestNewProtocolFromStringUDP(t *testing.T) {
	testNewProtocolFromString("UDP", t)
}

func TestNewProtocolFromStringTCP(t *testing.T) {
	testNewProtocolFromString("TCP", t)
}

func TestNewProtocolFromStringICMP(t *testing.T) {
	testNewProtocolFromString("ICMP", t)
}

func TestNewProtocolFromStringGracefullyHandlesCapitalisation(t *testing.T) {
	protocol, _ := NewProtocolFromString("iCmP")
	if protocol.String() != "ICMP" {
		t.Errorf("Protocol returned the wrong string value, got %v want %v", protocol.String(), "ICMP")
	}
}

func TestNewProtocolFromInvalidString(t *testing.T) {
	_, err := NewProtocolFromString("This_value_sucks")
	if err != ErrInvalidProtocolString {
		t.Errorf("Protocol constructor did not return the expected error")
	}
}
