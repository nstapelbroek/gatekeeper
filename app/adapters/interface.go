package adapters

import (
	"github.com/nstapelbroek/gatekeeper/domain"
)

// Adapter interface is used as an interface for the actual adapter implementations
type Adapter interface {
	CreateRule(rule domain.Rule) (err error)
	DeleteRule(rule domain.Rule) (err error)
}
