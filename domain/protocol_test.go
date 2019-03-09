package domain

import (
	"testing"
)

func testNewProtocolFromString(testValue string, expectedValue string, t *testing.T) {
	protocol, _ := NewProtocolFromString(testValue)
	if protocol.String() != expectedValue {
		t.Errorf("Protocol returned the wrong string value, got %v want %v", protocol.String(), expectedValue)
	}
}

func TestNewProtocolFromStringUDP(t *testing.T) {
	testNewProtocolFromString("UDP", "udp", t)
	testNewProtocolFromString("uDp", "udp", t)
	testNewProtocolFromString("udp", "udp", t)
}

func TestNewProtocolFromStringTCP(t *testing.T) {
	testNewProtocolFromString("TCP", "tcp", t)
	testNewProtocolFromString("TcP", "tcp", t)
	testNewProtocolFromString("tcp", "tcp", t)
}

func TestNewProtocolFromStringICMP(t *testing.T) {
	testNewProtocolFromString("ICMP", "icmp", t)
	testNewProtocolFromString("iCmP", "icmp", t)
	testNewProtocolFromString("icmp", "icmp", t)
}

func TestNewProtocolFromInvalidString(t *testing.T) {
	_, err := NewProtocolFromString("This_value_is_invalid")
	if err != ErrInvalidProtocolString {
		t.Errorf("Protocol constructor did not return the expected error")
	}
}
