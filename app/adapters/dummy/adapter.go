package dummy

import (
	"github.com/nstapelbroek/gatekeeper/domain"
)

type adapter struct{}

func (a *adapter) ToString() string {
	return "dummy"
}

func (adapter *adapter) CreateRules(rules []domain.Rule) domain.AdapterResult {
	return domain.AdapterResult{
		Error:  nil,
	}
}

func (adapter *adapter) DeleteRules(rules []domain.Rule) domain.AdapterResult {
	return domain.AdapterResult{
		Error:  nil,
	}
}
