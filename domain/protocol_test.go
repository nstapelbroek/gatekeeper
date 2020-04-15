package domain

import (
	"github.com/stretchr/testify/assert"
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
	testNewProtocolFromString("17", "udp", t)
}

func TestNewProtocolFromStringTCP(t *testing.T) {
	testNewProtocolFromString("TCP", "tcp", t)
	testNewProtocolFromString("TcP", "tcp", t)
	testNewProtocolFromString("tcp", "tcp", t)
	testNewProtocolFromString("6", "tcp", t)
}

func TestNewProtocolFromStringICMP(t *testing.T) {
	testNewProtocolFromString("ICMP", "icmp", t)
	testNewProtocolFromString("iCmP", "icmp", t)
	testNewProtocolFromString("icmp", "icmp", t)
	testNewProtocolFromString("1", "icmp", t)
}

func TestNewProtocolFromInvalidString(t *testing.T) {
	_, err := NewProtocolFromString("This_value_is_invalid")
	if err != ErrInvalidProtocolString {
		t.Errorf("Protocol constructor did not return the expected error")
	}
}

func TestTCPProtocolToIANANumber(t *testing.T) {
	p := Protocol(TCP)

	assert.Equal(t, 6, p.ProtocolNumber())
}

func TestUDPProtocolToIANANumber(t *testing.T) {
	p := Protocol(UDP)

	assert.Equal(t, 17, p.ProtocolNumber())
}

func TestICMProtocolToIANANumber(t *testing.T) {
	p := Protocol(ICMP)

	assert.Equal(t, 1, p.ProtocolNumber())
}
