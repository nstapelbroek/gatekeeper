package adapters

import (
	"github.com/spf13/viper"
	"github.com/nstapelbroek/gatekeeper/adapters/vultr"
)

type AdapterFactory struct {
	config *viper.Viper
}

func NewAdapterFactory(config *viper.Viper) *AdapterFactory {
	f := new(AdapterFactory)
	f.config = config

	return f
}

func (c AdapterFactory) GetAdapter() (a Adapter) {
	// currently, the only adapter implemented is Vultr so we'll return that one
	a = vultr.Adapter{
		ApiKey:          c.config.GetString("vultr_api_key"),
		FireWallGroupId: c.config.GetString("vultr_firewall_group_id"),
	}

	return
}
