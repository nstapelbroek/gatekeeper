package vultr

type RuleListResponse map[string]Rule

type Rule struct {
	RuleNumber int    `json:rulenumber`
	Action     string `json:action`
	Protocol   string `json:protocol`
	Port       string `json:port`
	Subnet     string `json:subnet`
	SubnetSize int    `json:subnet_size`
}