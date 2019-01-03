// Package adapters holds the generic struct's and interfaces for the adapter implementations and resolvers
package adapters

import (
	"github.com/nstapelbroek/gatekeeper/adapters/vultr"
	"github.com/spf13/viper"
)

// AdapterFactory forms a resolver that is responsible for transforming your configuration into adapter instances
type AdapterFactory struct {
	config *viper.Viper
}

// NewAdapterFactory is the constructor for the AdapterFactory
func NewAdapterFactory(config *viper.Viper) *AdapterFactory {
	f := new(AdapterFactory)
	f.config = config

	return f
}

// GetAdapter will return a adapter implementation based on your environment setup
func (c AdapterFactory) GetAdapter() (a Adapter) {
	// currently, the only adapter implemented is Vultr so we'll return that one
	a, err := vultr.NewVultrAdapter(
		c.config.GetString("vultr_api_key"),
		c.config.GetString("vultr_firewall_group"),
	)

	if err != nil {
		panic(err.Error())
	}

	return
}
