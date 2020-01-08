package vpc

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"strconv"
	"strings"
)

// ACLEntryCollection is a collection object for ACLEntryRuleNumbers
type ACLEntryCollection struct {
	rules map[string]int64
}

// NewACLEntryCollection is a constructor for ACLEntryCollection
func NewACLEntryCollection(entries []ec2.NetworkAclEntry) *ACLEntryCollection {
	c := &ACLEntryCollection{rules: make(map[string]int64)}
	for _, aclEntry := range entries {
		if *aclEntry.Egress || aclEntry.RuleAction != ec2.RuleActionAllow {
			continue
		}
		c.rules[c.aclEntryToUniqueKey(aclEntry)] = *aclEntry.RuleNumber
	}

	return c
}

// FindACLRuleNumberByRule will map a domain.rule to an ACLEntryRuleNumber in the internal collection
func (c *ACLEntryCollection) FindACLRuleNumberByRule(rule domain.Rule) *int64 {
	ruleNumber, exists := c.rules[c.ruleToUniqueKey(rule)]
	if exists {
		return &ruleNumber
	}
	return nil
}

func (c *ACLEntryCollection) ruleToUniqueKey(rule domain.Rule) string {
	return strings.Join([]string{
		rule.IPNet.String(),
		strconv.Itoa(rule.Protocol.ProtocolNumber()),
		strconv.Itoa(rule.Port.BeginPort),
		strconv.Itoa(rule.Port.EndPort),
	}, "-")
}

func (c *ACLEntryCollection) aclEntryToUniqueKey(entry ec2.NetworkAclEntry) string {
	cidr := entry.Ipv6CidrBlock
	if cidr == nil {
		cidr = entry.CidrBlock
	}

	return strings.Join([]string{
		*cidr,
		*entry.Protocol,
		strconv.FormatInt(*entry.PortRange.From, 10),
		strconv.FormatInt(*entry.PortRange.To, 10)},
		"-",
	)
}
