package adapters

import (
	"errors"
	"github.com/nstapelbroek/gatekeeper/domain"
	"sync"
)

type AdapterDispatcher struct {
	adapterInstances []domain.Adapter
}

type DispatchResult struct {
	SuccessfulDispatches []domain.AdapterResult
	FailedDispatches     []domain.AdapterResult
}

func (dr DispatchResult) HasFailures() bool {
	return len(dr.FailedDispatches) > 0
}

func NewAdapterDispatcher(adapterInstances []domain.Adapter) (*AdapterDispatcher, error) {
	if len(adapterInstances) == 0 {
		return nil, errors.New("no adapters configured")
	}

	d := AdapterDispatcher{
		adapterInstances: adapterInstances,
	}

	return &d, nil
}

func (ad AdapterDispatcher) Open(rules []domain.Rule) DispatchResult {
	return ad.dispatch(rules, "create")
}

func (ad AdapterDispatcher) Close(rules []domain.Rule) DispatchResult {
	return ad.dispatch(rules, "delete")
}

func (ad AdapterDispatcher) dispatch(rules []domain.Rule, action string) (dispatchResult DispatchResult) {
	resultChannel := make(chan domain.AdapterResult)
	var wg sync.WaitGroup

	for _, adapter := range ad.adapterInstances {
		wg.Add(1)
		go func(a domain.Adapter) {
			defer wg.Done()
			if action == "create" {
				resultChannel <- a.CreateRules(rules)
			} else if action == "delete" {
				resultChannel <- a.DeleteRules(rules)
			}
		}(adapter)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for result := range resultChannel {
		if result.IsSuccessful() {
			dispatchResult.SuccessfulDispatches = append(dispatchResult.SuccessfulDispatches, result)
			continue
		}

		dispatchResult.FailedDispatches = append(dispatchResult.FailedDispatches, result)
	}

	return dispatchResult
}
