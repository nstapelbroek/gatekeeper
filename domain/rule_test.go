package domain

import (
	"net"
	"testing"
)

func TestRuleCanConvertToString(t *testing.T) {
	direction, _ := NewDirectionFromString("outbound")
	rule := Rule{
		Direction: direction,
		Protocol:  Protocol(1),
		IPNet:     net.IPNet{IP: net.IP{192, 168, 1, 12}, Mask: net.IPMask{255, 255, 255, 0}},
		Port:      PortRange{20, 22},
	}

	println(rule.String())
	if rule.String() != "Rule-outbound-tcp-192.168.1.12-ffffff00-20-22" {
		t.Errorf("Rule converstion to string did not match expected result")
	}
}
