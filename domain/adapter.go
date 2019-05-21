package domain

type AdapterResult struct {
	Output string
	Error  error
}

func (r AdapterResult) IsSuccessful() bool {
	return r.Error == nil
}

type Adapter interface {
	ToString() string
	CreateRules(rule []Rule) AdapterResult
	DeleteRules(rule []Rule) AdapterResult
}
