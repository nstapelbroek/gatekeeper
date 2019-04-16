// Package adapters holds the generic structs and interfaces for the adapter implementations and resolvers
package adapters

import (
	"github.com/nstapelbroek/gatekeeper/app/adapters/digitalocean"
	"github.com/nstapelbroek/gatekeeper/app/adapters/vultr"
	"github.com/nstapelbroek/gatekeeper/domain"
	"github.com/spf13/viper"
)

type AdapterFactory struct {
	config *viper.Viper
}

func NewAdapterFactory(config *viper.Viper) *AdapterFactory {
	f := new(AdapterFactory)
	f.config = config

	return f
}

func (c AdapterFactory) GetAdapters() (adapterCollection []domain.Adapter) {
	config := c.config

	doToken := config.GetString("digitalocean_personal_access_token")
	doFirewallId := config.GetString("digitalocean_firewall_id")
	if len(doToken) > 0 && len(doFirewallId) > 0 {
		adapterCollection = append(adapterCollection, digitalocean.NewDigitalOceanAdapter(doToken, doFirewallId))
	}

	vultrToken := config.GetString("vultr_personal_access_token")
	vultrFirewallId := config.GetString("vultr_firewall_id")
	if len(vultrToken) > 0 && len(vultrFirewallId) > 0 {
		adapterCollection = append(adapterCollection, vultr.NewVultrAdapter(vultrToken, vultrFirewallId))
	}

	return
}
