package digitalocean

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/nstapelbroek/gatekeeper/domain"
	"golang.org/x/oauth2"
)

// Adapter is a DigitalOcean API implementation of the domain.Adapter interface
type Adapter struct {
	client     *godo.Client
	accessToken string
	firewallID string
}

// NewDigitalOceanAdapter is a constructor for Adapter
func NewDigitalOceanAdapter(personalAccessToken string, firewallIdentifier string) *Adapter {
	adapter := new(Adapter)
	adapter.accessToken = personalAccessToken
	adapter.firewallID = firewallIdentifier

	oauthClient := oauth2.NewClient(context.Background(), adapter)
	doClient := godo.NewClient(oauthClient)
	adapter.client = doClient

	return adapter
}

// Token is used to provide the DigitalOcean API client with an access token
func (a *Adapter) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: a.accessToken,
	}

	return token, nil
}

func (a *Adapter) ToString() string {
	return "digitalocean"
}

func (a *Adapter) newRequestFromDomainRule(rules []domain.Rule) *godo.FirewallRulesRequest {
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

func (a *Adapter) CreateRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.AddRules(context.TODO(), a.firewallID, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}

func (a *Adapter) DeleteRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.RemoveRules(context.TODO(), a.firewallID, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}
