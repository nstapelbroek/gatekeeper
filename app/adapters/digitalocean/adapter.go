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

func (a *adapter) newRequestFromDomainRule(rule domain.Rule) *godo.FirewallRulesRequest {
	doSource := &godo.Sources{Addresses: []string{rule.IPNet.String()}}
	rulesRequest := &godo.FirewallRulesRequest{
		InboundRules: []godo.InboundRule{{
			Protocol:  rule.Protocol.String(),
			PortRange: rule.Port.String(),
			Sources:   doSource,
		}},
	}

	return rulesRequest
}
func (a *adapter) CreateRule(rule domain.Rule) (err error) {
	_, err = a.client.Firewalls.AddRules(context.TODO(), a.firewallId, a.newRequestFromDomainRule(rule))
	return
}

func (a *adapter) DeleteRule(rule domain.Rule) (err error) {
	_, err = a.client.Firewalls.RemoveRules(context.TODO(), a.firewallId, a.newRequestFromDomainRule(rule))
	return
}
