package vultr

import (
	"github.com/JamesClonk/vultr/lib"
	"github.com/nstapelbroek/gatekeeper/domain"
)

type adapter struct {
	client           *lib.Client
	firewallGroupId  string
	ruleNumbersIndex map[string]int
}

func NewVultrAdapter(apiKey string, firewallGroupID string) *adapter {
	vultrClient := lib.NewClient(apiKey, nil)

	adapter := new(adapter)
	adapter.client = vultrClient
	adapter.firewallGroupId = firewallGroupID
	adapter.ruleNumbersIndex = make(map[string]int)

	return adapter
}

func (a *adapter) CreateRule(rule domain.Rule) (err error) {
	_, keyExists := a.ruleNumbersIndex[rule.String()]
	if keyExists {
		return // Block subsequent rule requests util it's removed by the timeout
	}

	ruleNumber, err := a.client.CreateFirewallRule(a.firewallGroupId, rule.Protocol.String(), rule.Port.String(), &rule.IPNet)
	if err == nil {
		a.ruleNumbersIndex[rule.String()] = ruleNumber
	}

	return
}

func (a *adapter) DeleteRule(rule domain.Rule) (err error) {
	ruleNumber, keyExists := a.ruleNumbersIndex[rule.String()]
	if !keyExists {
		return
	}

	delete(a.ruleNumbersIndex, rule.String())
	return a.client.DeleteFirewallRule(ruleNumber, a.firewallGroupId)
}
