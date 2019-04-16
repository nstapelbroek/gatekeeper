package dummy

import (
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"testing"
)

func getRules() []domain.Rule {
	protocol, _ := domain.NewProtocolFromString("TCP")
	direction, _ := domain.NewDirectionFromString("inbound")

	r1 := domain.Rule{
		IPNet:     net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(32, 32)},
		Port:      domain.NewSinglePort(22),
		Protocol:  protocol,
		Direction: direction,
	}

	r2 := domain.Rule{
		IPNet:     net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(32, 32)},
		Port:      domain.NewSinglePort(20),
		Protocol:  protocol,
		Direction: direction,
	}

	return []domain.Rule{r1, r2}
}

func TestAdapter_CreateRule(t *testing.T) {
	adapterInstance := adapter{}
	if !adapterInstance.CreateRules(getRules()).IsSuccessful() {
		t.Error("Dummy adapter is not supposed to do anything but it somehow created an error")
	}
}

func TestAdapter_DeleteRule(t *testing.T) {
	adapterInstance := adapter{}
	if !adapterInstance.DeleteRules(getRules()).IsSuccessful() {
		t.Error("Dummy adapter is not supposed to do anything but it somehow created an error")
	}
}
