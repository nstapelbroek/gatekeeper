package vpc

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"strconv"
	"strings"
)

type aclEntryCollection struct {
	rules map[string]int64
}

func NewACLEntryCollection(entries []ec2.NetworkAclEntry) *aclEntryCollection {
	c := &aclEntryCollection{rules: make(map[string]int64)}
	for _, aclEntry := range entries {
		if *aclEntry.Egress == true || aclEntry.RuleAction != ec2.RuleActionAllow {
			continue
		}
		c.rules[c.aclEntryToUniqueKey(aclEntry)] = *aclEntry.RuleNumber
	}

	return c
}

func (c *aclEntryCollection) FindAclRuleNumberByRule(rule domain.Rule) *int64 {
	ruleNumber, exists := c.rules[c.ruleToUniqueKey(rule)]
	if exists {
		return &ruleNumber
	}
	return nil
}

func (c *aclEntryCollection) ruleToUniqueKey(rule domain.Rule) string {
	return strings.Join([]string{
		rule.IPNet.String(),
		strconv.Itoa(rule.Protocol.ProtocolNumber()),
		strconv.Itoa(rule.Port.BeginPort),
		strconv.Itoa(rule.Port.EndPort),
	}, "-")
}

func (c *aclEntryCollection) aclEntryToUniqueKey(entry ec2.NetworkAclEntry) string {
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
