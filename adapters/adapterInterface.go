package adapters

import (
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
)

// Adapter interface is used as an interface for the actual adapter implementations
type Adapter interface {
	CreateRule(rule firewall.Rule) (err error)
	DeleteRule(rule firewall.Rule) (err error)
}
