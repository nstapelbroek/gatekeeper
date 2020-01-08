package digitalocean

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/nstapelbroek/gatekeeper/domain"
	"golang.org/x/oauth2"
)

// Adapter is a DigitalOcean API implementation of the domain.Adapter interface
type Adapter struct {
	client      *godo.Client
	tokenSource tokenSource
	firewallID  string
}

type tokenSource struct {
	AccessToken string
}

// Token is used to provide the DigitalOcean API client with an access token
func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// NewDigitalOceanAdapter is a constructor for Adapter
func NewDigitalOceanAdapter(personalAccessToken string, firewallIdentifier string) *Adapter {
	a := new(Adapter)
	a.tokenSource = tokenSource{AccessToken: personalAccessToken}
	a.firewallID = firewallIdentifier

	oauthClient := oauth2.NewClient(context.Background(), &a.tokenSource)
	doClient := godo.NewClient(oauthClient)
	a.client = doClient

	return a
}

// ToString satisfies the domain.Adapter interface
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

// CreateRules satisfies the domain.Adapter interface
func (a *Adapter) CreateRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.AddRules(context.TODO(), a.firewallID, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}

// DeleteRules satisfies the domain.Adapter interface
func (a *Adapter) DeleteRules(rules []domain.Rule) domain.AdapterResult {
	_, err := a.client.Firewalls.RemoveRules(context.TODO(), a.firewallID, a.newRequestFromDomainRule(rules))

	return domain.AdapterResult{Error: err}
}
