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

type adapter struct {
	client           *ec2.Client
	networkAclId     string
	allowedRuleRange aclRuleNumberRange
}

type aclRuleNumberRange struct {
	min int64
	max int64
}

func NewAWSNetworkACLAdapter(client *ec2.Client, networkAclId string, numberRange string) *adapter {
	nRange := strings.SplitN(numberRange, "-", 2)
	min, _ := strconv.ParseInt(nRange[0], 10, 0)
	max, _ := strconv.ParseInt(nRange[1], 10, 0)

	return &adapter{
		client:           client,
		networkAclId:     networkAclId,
		allowedRuleRange: aclRuleNumberRange{min: min, max: max},
	}
}

func (a *adapter) ToString() string {
	return "aws-network-acl"
}

func (a *adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	currentEntries := a.getPersistedAclEntries()
	availableRuleNumbers, err := a.calculateAvailableRuleNumbers(currentEntries, len(rules))
	if err != nil {
		return domain.AdapterResult{Error: err}
	}

	for i, rule := range rules {
		if currentEntries.FindAclRuleNumberByRule(rule) != nil {
			return domain.AdapterResult{Error: errors.New("rule is already set")}
		}

		input := ec2.CreateNetworkAclEntryInput{
			CidrBlock:    aws.String(rule.IPNet.String()),
			Egress:       aws.Bool(rule.Direction.IsOutbound()),
			NetworkAclId: aws.String(a.networkAclId),
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

func (a *adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	currentEntries := a.getPersistedAclEntries()
	for _, rule := range rules {
		ruleNumber := currentEntries.FindAclRuleNumberByRule(rule)
		if ruleNumber == nil {
			// todo log that a rule could not be found and is therefore ignored in the cleanup
			continue
		}

		input := ec2.DeleteNetworkAclEntryInput{
			Egress:       aws.Bool(rule.Direction.IsOutbound()),
			NetworkAclId: aws.String(a.networkAclId),
			RuleNumber:   ruleNumber,
		}

		req := a.client.DeleteNetworkAclEntryRequest(&input)
		_, _ = req.Send(context.TODO()) // TODO error logging when it's a background task.
	}

	return domain.AdapterResult{}
}

func (a *adapter) getPersistedAclEntries() *aclEntryCollection {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{
			a.networkAclId,
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

func (a *adapter) calculateAvailableRuleNumbers(entries *aclEntryCollection, requestedCount int) ([]int64, error) {
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
