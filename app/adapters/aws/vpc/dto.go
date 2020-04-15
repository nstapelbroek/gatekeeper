package vpc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"net"
	"strconv"
)

type EntryCollection struct {
	entries       map[string][]ec2.NetworkAclEntry
	highestNumber int64
}

func (c EntryCollection) findByRule(rule domain.Rule) *ec2.NetworkAclEntry {
	persistedRules, found := c.entries[rule.String()]
	if !found {
		return nil
	}

	for i := range persistedRules {
		if persistedRules[i].RuleAction == ec2.RuleActionAllow {
			return &persistedRules[i]
		}
	}

	return nil
}

// NewACLEntryCollection is a constructor for ACLEntryCollection
func NewEntryCollection(entries []ec2.NetworkAclEntry) *EntryCollection {
	c := &EntryCollection{
		entries:       make(map[string][]ec2.NetworkAclEntry),
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

func networkEntryToRule(entry ec2.NetworkAclEntry) (*domain.Rule, error) {
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
			BeginPort: *entry.PortRange.From,
			EndPort:   *entry.PortRange.To,
		},
	}, nil
}

func createAddEntryInput(rule domain.Rule, networkAclID string, ruleNumber int64) *ec2.CreateNetworkAclEntryInput {
	input := &ec2.CreateNetworkAclEntryInput{
		CidrBlock:    aws.String(rule.IPNet.String()),
		Egress:       aws.Bool(rule.Direction.IsOutbound()),
		NetworkAclId: aws.String(networkAclID),
		PortRange:    &ec2.PortRange{From: aws.Int64(rule.Port.BeginPort), To: aws.Int64(rule.Port.EndPort)},
		Protocol:     aws.String(strconv.Itoa(rule.Protocol.ProtocolNumber())),
		RuleAction:   "allow",
		RuleNumber:   aws.Int64(ruleNumber),
	}

	if rule.IPNet.IP.To4() == nil {
		input.Ipv6CidrBlock = input.CidrBlock
		input.CidrBlock = nil
	}

	return input
}

func createDeleteEntryInput(persistedRule *ec2.NetworkAclEntry, networkAclID string) *ec2.DeleteNetworkAclEntryInput {
	return &ec2.DeleteNetworkAclEntryInput{
		Egress:       persistedRule.Egress,
		NetworkAclId: aws.String(networkAclID),
		RuleNumber:   persistedRule.RuleNumber,
	}
}
