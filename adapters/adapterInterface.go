package adapters

import (
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
)

type Adapter interface {
	CreateRule(rule firewall.Rule) (err error)
	DeleteRule(rule firewall.Rule) (err error)
}
