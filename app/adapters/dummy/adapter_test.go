package dummy

import (
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"testing"
)

func getRule() domain.Rule {
	protocol, _ := firewall.NewProtocolFromString("TCP")
	direction, _ := firewall.NewDirectionFromString("inbound")

	r := domain.Rule{
		IPNet:     net.ParseIP("127.0.0.1"),
		Port:      firewall.NewSinglePort(22),
		Protocol:  protocol,
		Direction: direction,
	}

	return r
}

func TestAdapter_CreateRule(t *testing.T) {
	adapterInstance := adapter{}
	if adapterInstance.CreateRule(getRule()) != nil {
		t.Error("Dummy adapter is not supposed to do anything but it somehow created an error")
	}
}

func TestAdapter_DeleteRule(t *testing.T) {
	adapterInstance := adapter{}
	if adapterInstance.DeleteRule(getRule()) != nil {
		t.Error("Dummy adapter is not supposed to do anything but it somehow created an error")
	}
}
