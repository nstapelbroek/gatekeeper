package dummy

import (
	"github.com/nstapelbroek/gatekeeper/domain"
)

type adapter struct{}

// NewDummyAdapter will create a new dummy adapter object for testing purposes.
func NewDummyAdapter() *adapter {
	a := new(adapter)
	return a
}

func (adapter *adapter) CreateRule(rule domain.Rule) (err error) {
	return
}

func (adapter *adapter) DeleteRule(rule domain.Rule) (err error) {
	return
}
