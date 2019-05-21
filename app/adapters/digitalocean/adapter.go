package digitalocean

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/nstapelbroek/gatekeeper/domain"
	"golang.org/x/oauth2"
)

type adapter struct {
	client      *godo.Client
	accessToken string
	firewallId  string
}

func NewDigitalOceanAdapter(personalAccessToken string, firewallIdentifier string) *adapter {
	adapter := new(adapter)
	adapter.accessToken = personalAccessToken
	adapter.firewallId = firewallIdentifier

	oauthClient := oauth2.NewClient(context.Background(), adapter)
	doClient := godo.NewClient(oauthClient)
	adapter.client = doClient

	return adapter
}

func (a *adapter) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: a.accessToken,
	}

	return token, nil
}

func (a *adapter) ToString() string {
	return "digitalocean"
}

func (a *adapter) newRequestFromDomainRule(rules []domain.Rule) *godo.FirewallRulesRequest {
	rulesRequest := &godo.FirewallRulesRequest{
		InboundRules: []godo.InboundRule{},
	}

	for _, rule := range rules {
		doSource := &godo.Sources{Addresses: []string{rule.IPNet.String()}}
		rulesRequest.InboundRules = append(rulesRequest.InboundRules,
			godo.InboundRule{
				Protocol:  rule.Protocol.String(),
				PortRange: rule.Port.String(),
				Sources:   doSource,
			},
		)
	}

	return rulesRequest
}

func (a *adapter) CreateRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.AddRules(context.TODO(), a.firewallId, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}

func (a *adapter) DeleteRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.RemoveRules(context.TODO(), a.firewallId, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}
