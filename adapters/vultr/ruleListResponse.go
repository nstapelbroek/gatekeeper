package vultr

// RuleListResponse is the response expected when sending a RuleListRequest
type RuleListResponse map[string]Rule

// Rule is a struct used when unmarshalling the RuleListResponse response's body
type Rule struct {
	RuleNumber int    `json:rulenumber`
	Action     string `json:action`
	Protocol   string `json:protocol`
	Port       string `json:port`
	Subnet     string `json:subnet`
	SubnetSize int    `json:subnet_size`
}