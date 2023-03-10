package adapters

import (
	"errors"
	"sync"

	"github.com/nstapelbroek/gatekeeper/domain"
	"go.uber.org/zap"
)

// AdapterDispatcher will coordinate calling all configured adapters in a WaitGroup
type AdapterDispatcher struct {
	adapterInstances *[]domain.Adapter
	logger           *zap.Logger
}

type wrappedAdapterResult struct {
	Name   string
	Result domain.AdapterResult
}

// NewAdapterDispatcher is a constructor method for AdapterDispatcher
func NewAdapterDispatcher(adapterInstances *[]domain.Adapter, logger *zap.Logger) (*AdapterDispatcher, error) {
	d := AdapterDispatcher{
		adapterInstances: adapterInstances,
		logger:           logger,
	}

	return &d, nil
}

// Open will call CreateRules on all configured adapters in the AdapterDispatcher
func (ad AdapterDispatcher) Open(rules []domain.Rule) (map[string]string, error) {
	return ad.dispatch(rules, "create")
}

// Close will call CreateRules on all configured adapters in the AdapterDispatcher
func (ad AdapterDispatcher) Close(rules []domain.Rule) (map[string]string, error) {
	return ad.dispatch(rules, "delete")
}

func (ad AdapterDispatcher) dispatch(rules []domain.Rule, action string) (map[string]string, error) {
	resultChannel := make(chan wrappedAdapterResult)
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

			resultChannel <- wrappedAdapterResult{
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

func (ad AdapterDispatcher) processDispatchResults(resultChannel chan wrappedAdapterResult) (dispatchResult map[string]string, err error) {
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
		ad.logger.Debug(
			"Result from API",
			zap.String("adapter", name),
			zap.String("output", output),
			zap.Bool("error", hasErr),
		)
	}

	if hasErr {
		err = errors.New("failed applying some rules")
	}

	return dispatchResult, err
}
