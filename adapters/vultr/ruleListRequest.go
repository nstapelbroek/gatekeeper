package vultr

import (
	"fmt"
	"net/http"
	"strings"
)

type RuleListRequest struct {
	request         *http.Request
	firewallGroupId string
	direction       string
	ipType          string
}

func NewRuleListRequest(ApiKey string, FirewallGroupId string, IPType string) *RuleListRequest {
	r := new(RuleListRequest)
	r.firewallGroupId = FirewallGroupId
	r.direction = "in"
	r.ipType = IPType

	var endpoint = "https://api.vultr.com/v1/firewall/rule_list?" + r.getQueryString()
	r.request, _ = http.NewRequest(http.MethodGet, endpoint, nil)
	r.request.Header.Set("Api-Key", ApiKey)

	return r
}

func (r RuleListRequest) getQueryString() string {
	return fmt.Sprintf(
		"FIREWALLGROUPID=%s&direction=%s&ip_type=%s",
		r.firewallGroupId,
		strings.ToLower(r.direction),
		strings.ToLower(r.ipType),
	)
}
