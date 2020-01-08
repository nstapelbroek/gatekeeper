package domain

// AdapterResult will act as a value object of wrapping more information about an Create or Delete execution later
type AdapterResult struct {
	Output string
	Error  error
}

// IsSuccessful is a method to generalize the way of checking an erroneous AdapterResult
func (r AdapterResult) IsSuccessful() bool {
	return r.Error == nil
}

// Adapter interface is the definition where all adapter implementations have to ad-here to
type Adapter interface {
	ToString() string
	CreateRules(rule []Rule) AdapterResult
	DeleteRules(rule []Rule) AdapterResult
}
