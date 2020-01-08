package domain

import (
	"net"
	"strings"
	"time"
)

// OpenCommand will act as a command object we can pass around in the app
type OpenCommand struct {
	Rules     []Rule
	Timeout   time.Duration
	IPAddress net.IPNet
}

// NewOpenCommand is a constructor for the OpenCommand
func NewOpenCommand(IP string, timeout int64, boilerplateRules []Rule) *OpenCommand {
	var o OpenCommand

	o.IPAddress = parseIP(IP)
	o.Timeout = parseTimeout(timeout)
	o.Rules = parseRules(o.IPAddress, boilerplateRules)

	return &o
}

func parseIP(stringIP string) net.IPNet {
	if strings.Contains(stringIP, "/") {
		ip, ipNet, _ := net.ParseCIDR(stringIP)
		return net.IPNet{IP: ip, Mask: ipNet.Mask}
	}

	ip := net.ParseIP(stringIP)
	if ip.To4() != nil {
		return net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
	}

	// Assuming IPv6 here because conversion to IPv4 failed
	return net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
}

func parseTimeout(rawTimeout int64) time.Duration {
	return time.Duration(rawTimeout) * time.Second
}

func parseRules(ipAddress net.IPNet, defaultRules []Rule) []Rule {
	rules := make([]Rule, len(defaultRules))
	copy(rules, defaultRules)
	for index := range rules {
		rules[index].IPNet = ipAddress
	}

	return rules
}
