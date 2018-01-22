package adapters

import (
	"github.com/spf13/viper"
	"github.com/nstapelbroek/gatekeeper/adapters/vultr"
)

type AdapterFactory struct {
	Config *viper.Viper
}

func (c AdapterFactory) GetAdapter() (a Adapter) {
	// currently, the only adapter implemented is Vultr so we'll return that one
	a = vultr.Adapter{
		ApiKey:          c.Config.GetString("vultr_api_key"),
		FireWallGroupId: c.Config.GetString("vultr_firewall_group_id"),
	}

	return
}
