package adapters

import (
	"errors"
	"fmt"
	"github.com/nstapelbroek/gatekeeper/domain"
)

type AdapterDispatcher struct {
	inputChannels []chan []domain.Rule
	resultChannel chan domain.AdapterResult
}

func adapterWorker(adapter domain.Adapter, input <-chan []domain.Rule, results chan<- domain.AdapterResult) {
	fmt.Sprintln("worker for adapter %s stated", &adapter)
	for rules := range input {
		fmt.Println("recieved job!")
		results <- adapter.CreateRules(rules)
		fmt.Println("finished job")
	}
}

func NewAdapterDispatcher(adapterInstances []domain.Adapter) (*AdapterDispatcher, error) {
	if len(adapterInstances) == 0 {
		return nil, errors.New("no adapters configured")
	}

	d := AdapterDispatcher{
		resultChannel: make(chan domain.AdapterResult),
	}

	for _, adapter := range adapterInstances {
		input := make(chan []domain.Rule)
		d.inputChannels = append(d.inputChannels, input)
		go adapterWorker(adapter, input, d.resultChannel)
	}

	return &d, nil
}

func (ad AdapterDispatcher) Open(rules []domain.Rule) (results []domain.AdapterResult) {
	for _, inputChannel := range ad.inputChannels {
		inputChannel <- rules
	}

	for len(results) != len(ad.inputChannels) {
		result := <-ad.resultChannel
		println(result.IsSuccessful())
		results = append(results, result)
	}

	return results
}

func (ad AdapterDispatcher) Close(rules []domain.Rule) (results []domain.AdapterResult) {
	fmt.Println("close is not implemented yet...")
	return
}
