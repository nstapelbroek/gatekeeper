package adapters

import (
	"errors"
	"github.com/nstapelbroek/gatekeeper/domain"
	"sync"
)

type AdapterDispatcher struct {
	adapterInstances *[]domain.Adapter
}

type WrappedAdapterResult struct {
	Name   string
	Result domain.AdapterResult
}

func NewAdapterDispatcher(adapterInstances *[]domain.Adapter) (*AdapterDispatcher, error) {
	d := AdapterDispatcher{
		adapterInstances: adapterInstances,
	}

	return &d, nil
}

func (ad AdapterDispatcher) Open(rules []domain.Rule) (map[string]string, error) {
	return ad.dispatch(rules, "create")
}

func (ad AdapterDispatcher) Close(rules []domain.Rule) (map[string]string, error) {
	return ad.dispatch(rules, "delete")
}

func (ad AdapterDispatcher) dispatch(rules []domain.Rule, action string) (map[string]string, error) {
	resultChannel := make(chan WrappedAdapterResult)
	var wg sync.WaitGroup

	for _, adapter := range *ad.adapterInstances {
		wg.Add(1)
		go func(a domain.Adapter) {
			defer wg.Done()

			var result domain.AdapterResult
			if action == "create" {
				result = a.CreateRules(rules)
			} else if action == "delete" {
				result = a.DeleteRules(rules)
			}

			resultChannel <- WrappedAdapterResult{
				Name:   a.ToString(),
				Result: result,
			}
		}(adapter)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	return ad.processDispatchResults(resultChannel)
}

func (ad AdapterDispatcher) processDispatchResults(resultChannel chan WrappedAdapterResult) (dispatchResult map[string]string, err error) {
	hasErr := false
	dispatchResult = make(map[string]string)

	for wrappedResult := range resultChannel {
		name := wrappedResult.Name
		result := &wrappedResult.Result
		output := result.Output

		if result.Error != nil {
			hasErr = true
			output = result.Error.Error()
		}

		dispatchResult[name] = output
	}

	if hasErr {
		err = errors.New("failed applying some rules")
	}

	return dispatchResult, err
}
