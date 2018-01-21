package vultr

import (
	"fmt"
	"net/http"
	"strings"
)

type RuleDeleteRequest struct {
	request         *http.Request
	firewallGroupId string
	ruleNumber      int
}

func NewRuleDeleteRequest(ApiKey string, FirewallGroupId string, RuleNumber int) *RuleDeleteRequest {
	r := new(RuleDeleteRequest)
	r.firewallGroupId = FirewallGroupId
	r.ruleNumber = RuleNumber

	requestBody := strings.NewReader(r.getBodyString())
	r.request, _ = http.NewRequest(http.MethodPost, "https://api.vultr.com/v1/firewall/rule_delete", requestBody)
	r.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.request.Header.Set("Api-Key", ApiKey)

	return r
}

func (r RuleDeleteRequest) getBodyString() string {
	return fmt.Sprintf(
		"FIREWALLGROUPID=%s&rulenumber=%d",
		r.firewallGroupId,
		r.ruleNumber,
	)
}
