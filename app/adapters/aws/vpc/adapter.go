package vpc

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"strconv"
	"strings"
)

// Adapter is a AWS VPC API implementation of the domain.Adapter interface
type Adapter struct {
	client           *ec2.Client
	networkACLID     string
	allowedRuleRange aclRuleNumberRange
}

type aclRuleNumberRange struct {
	min int64
	max int64
}

// NewAWSNetworkACLAdapter is a constructor for Adapter
func NewAWSNetworkACLAdapter(client *ec2.Client, networkACLID string, numberRange string) *Adapter {
	nRange := strings.SplitN(numberRange, "-", 2)
	min, _ := strconv.ParseInt(nRange[0], 10, 0)
	max, _ := strconv.ParseInt(nRange[1], 10, 0)

	return &Adapter{
		client:           client,
		networkACLID:     networkACLID,
		allowedRuleRange: aclRuleNumberRange{min: min, max: max},
	}
}

// ToString satisfies the domain.Adapter interface
func (a *Adapter) ToString() string {
	return "aws-network-acl"
}

// CreateRules satisfies the domain.Adapter interface
func (a *Adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	currentEntries := a.getPersistedACLEntries()
	availableRuleNumbers, err := a.calculateAvailableRuleNumbers(currentEntries, len(rules))
	if err != nil {
		return domain.AdapterResult{Error: err}
	}

	for i, rule := range rules {
		if currentEntries.FindACLRuleNumberByRule(rule) != nil {
			return domain.AdapterResult{Error: errors.New("rule is already set")}
		}

		input := ec2.CreateNetworkAclEntryInput{
			CidrBlock:    aws.String(rule.IPNet.String()),
			Egress:       aws.Bool(rule.Direction.IsOutbound()),
			NetworkAclId: aws.String(a.networkACLID),
			PortRange:    &ec2.PortRange{From: aws.Int64(int64(rule.Port.BeginPort)), To: aws.Int64(int64(rule.Port.EndPort))},
			Protocol:     aws.String(strconv.Itoa(rule.Protocol.ProtocolNumber())),
			RuleAction:   "allow",
			RuleNumber:   &availableRuleNumbers[i],
		}

		if rule.IPNet.IP.To4() == nil {
			input.Ipv6CidrBlock = input.CidrBlock
			input.CidrBlock = nil
		}

		req := a.client.CreateNetworkAclEntryRequest(&input)
		_, err = req.Send(context.TODO())
		if err != nil {
			return domain.AdapterResult{Error: err}
		}
	}

	return domain.AdapterResult{}
}

// DeleteRules satisfies the domain.Adapter interface
func (a *Adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	currentEntries := a.getPersistedACLEntries()
	for _, rule := range rules {
		ruleNumber := currentEntries.FindACLRuleNumberByRule(rule)
		if ruleNumber == nil {
			// todo log that a rule could not be found and is therefore ignored in the cleanup
			continue
		}

		input := ec2.DeleteNetworkAclEntryInput{
			Egress:       aws.Bool(rule.Direction.IsOutbound()),
			NetworkAclId: aws.String(a.networkACLID),
			RuleNumber:   ruleNumber,
		}

		req := a.client.DeleteNetworkAclEntryRequest(&input)
		_, _ = req.Send(context.TODO()) // TODO error logging when it's a background task.
	}

	return domain.AdapterResult{}
}

func (a *Adapter) getPersistedACLEntries() *ACLEntryCollection {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{
			a.networkACLID,
		},
		Filters: []ec2.Filter{
			{
				Name:   aws.String("entry.rule-action"),
				Values: []string{"allow"},
			},
		},
	}

	req := a.client.DescribeNetworkAclsRequest(input)
	resp, err := req.Send(context.Background())
	if err != nil || len(resp.NetworkAcls) == 0 {
		return NewACLEntryCollection(nil) // todo log error
	}

	return NewACLEntryCollection(resp.NetworkAcls[0].Entries) // Assume only 1 result because we filtered
}

func (a *Adapter) calculateAvailableRuleNumbers(entries *ACLEntryCollection, requestedCount int) ([]int64, error) {
	takenNumbers := make(map[int64]bool)
	var availableNumbers []int64
	for _, ruleNumber := range entries.rules {
		if ruleNumber >= a.allowedRuleRange.min && ruleNumber <= a.allowedRuleRange.max {
			takenNumbers[ruleNumber] = true
		}
	}

	for availableNumber := a.allowedRuleRange.min; availableNumber < a.allowedRuleRange.max; availableNumber++ {
		if len(availableNumbers) == requestedCount {
			return availableNumbers, nil
		}

		if _, ok := takenNumbers[availableNumber]; !ok {
			availableNumbers = append(availableNumbers, availableNumber)
		}
	}

	return availableNumbers, errors.New("ran out of acl rule numbers to allocate")
}
