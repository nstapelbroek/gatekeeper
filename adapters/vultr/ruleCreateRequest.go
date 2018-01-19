package vultr

import (
	"fmt"
	"net/http"
	"strings"
)

type RuleCreateRequest struct {
	request         *http.Request
	firewallGroupId string
	direction       string
	ipType          string
	protocol        string
	subnet          string
	subnetSize      string
	port            string
}

func NewRuleCreateRequest(ApiKey string, FirewallGroupId string, IPType string, Protocol string, Subnet string, SubnetSize string, Port string) *RuleCreateRequest {
	r := new(RuleCreateRequest)
	r.firewallGroupId = FirewallGroupId
	r.direction = "in"
	r.ipType = IPType
	r.protocol = Protocol
	r.subnet = Subnet
	r.subnetSize = SubnetSize
	r.port = Port

	requestBody := strings.NewReader(r.getBodyString())
	r.request, _ = http.NewRequest(http.MethodPost, "https://api.vultr.com/v1/firewall/rule_create", requestBody)
	r.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.request.Header.Set("Api-Key", ApiKey)

	return r
}

func (r RuleCreateRequest) getBodyString() string {
	return fmt.Sprintf(
		"FIREWALLGROUPID=%s&direction=in&ip_type=%s&protocol=%s&subnet_size=%s&subnet=%s&port=%s",
		r.firewallGroupId,
		strings.ToLower(r.ipType),
		strings.ToLower(r.protocol),
		r.subnetSize,
		r.subnet,
		r.port,
	)
}
