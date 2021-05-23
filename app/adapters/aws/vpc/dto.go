package vpc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"strconv"
)

type EntryCollection struct {
	entries       map[string][]types.NetworkAclEntry
	highestNumber int32
}

func (c EntryCollection) findByRule(rule domain.Rule) (types.NetworkAclEntry, bool) {
	var foundRule types.NetworkAclEntry
	persistedRules, found := c.entries[rule.String()]
	if !found {
		return foundRule, false
	}

	for i := range persistedRules {
		if persistedRules[i].RuleAction == types.RuleActionAllow {
			return persistedRules[i], true
		}
	}

	return foundRule, false
}

// NewACLEntryCollection is a constructor for ACLEntryCollection
func NewEntryCollection(entries []types.NetworkAclEntry) *EntryCollection {
	c := &EntryCollection{
		entries:       make(map[string][]types.NetworkAclEntry),
		highestNumber: 0,
	}

	for _, entry := range entries {
		rule, err := networkEntryToRule(entry)
		if err != nil {
			// todo warn ignored rule
			continue
		}

		c.entries[rule.String()] = append(c.entries[rule.String()], entry)
		ruleNumber := *entry.RuleNumber
		if ruleNumber > c.highestNumber {
			c.highestNumber = ruleNumber
		}
	}

	return c
}

func networkEntryToRule(entry types.NetworkAclEntry) (*domain.Rule, error) {
	cidr := entry.Ipv6CidrBlock
	if cidr == nil {
		cidr = entry.CidrBlock
	}
	ip, ipNet, err := net.ParseCIDR(*cidr)
	if err != nil {
		return nil, err
	}

	protocol, err := domain.NewProtocolFromString(*entry.Protocol)
	if err != nil {
		return nil, err
	}

	direction := domain.Inbound
	if *entry.Egress {
		direction = domain.Outbound
	}

	return &domain.Rule{
		Direction: direction,
		Protocol:  protocol,
		IPNet:     net.IPNet{IP: ip, Mask: ipNet.Mask},
		Port: domain.PortRange{
			BeginPort: int64(*entry.PortRange.From),
			EndPort:   int64(*entry.PortRange.To),
		},
	}, nil
}

func createAddEntryInput(rule domain.Rule, networkAclID string, ruleNumber int32) *ec2.CreateNetworkAclEntryInput {
	input := &ec2.CreateNetworkAclEntryInput{
		CidrBlock:    aws.String(rule.IPNet.String()),
		Egress:       aws.Bool(rule.Direction.IsOutbound()),
		NetworkAclId: aws.String(networkAclID),
		PortRange:    &types.PortRange{From: aws.Int32(int32(rule.Port.BeginPort)), To: aws.Int32(int32(rule.Port.EndPort))},
		Protocol:     aws.String(strconv.Itoa(rule.Protocol.ProtocolNumber())),
		RuleAction:   "allow",
		RuleNumber:   aws.Int32(ruleNumber),
	}

	if rule.IPNet.IP.To4() == nil {
		input.Ipv6CidrBlock = input.CidrBlock
		input.CidrBlock = nil
	}

	return input
}

func createDeleteEntryInput(persistedRule types.NetworkAclEntry, networkAclID string) *ec2.DeleteNetworkAclEntryInput {
	return &ec2.DeleteNetworkAclEntryInput{
		Egress:       persistedRule.Egress,
		NetworkAclId: aws.String(networkAclID),
		RuleNumber:   persistedRule.RuleNumber,
	}
}
