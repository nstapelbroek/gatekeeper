package domain

import (
	"net"
	"strings"
	"time"
)

type OpenCommand struct {
	Rules     []Rule
	Timeout   time.Duration
	IpAddress net.IPNet
}

func NewOpenCommand(ip string, timeout int64, boilerplateRules []Rule) *OpenCommand {
	var o OpenCommand

	o.IpAddress = parseIp(ip)
	o.Timeout = parseTimeout(timeout)
	o.Rules = parseRules(o.IpAddress, boilerplateRules)

	return &o
}

func parseIp(stringIp string) net.IPNet {
	if strings.Contains(stringIp, "/") {
		ip, ipNet, _ := net.ParseCIDR(stringIp)
		return net.IPNet{IP: ip, Mask: ipNet.Mask}
	}

	ip := net.ParseIP(stringIp)
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
