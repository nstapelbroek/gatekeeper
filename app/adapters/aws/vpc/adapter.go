package vpc

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
)

type adapter struct {
	client          *ec2.Client
	NetworkAclId    string
	startRuleNumber int64
	ruleStepsTaken  int
	ruleStepSize    int
}

func NewAWSNetworkACLAdapter(client *ec2.Client, networkAclId string, startRuleNumber int64) *adapter {
	return &adapter{
		client:          client,
		NetworkAclId:    networkAclId,
		startRuleNumber: startRuleNumber,
		ruleStepsTaken:  0,
		ruleStepSize:    10,
	}
}

func (a *adapter) ToString() string {
	return "aws-network-acl"
}

func (a *adapter) getNextRuleNumber() *int64 {
	a.ruleStepsTaken = a.ruleStepsTaken + 1
	return aws.Int64(a.startRuleNumber + int64(a.ruleStepSize*a.ruleStepsTaken))
}

func (a *adapter) getProtocolNumber(protocol domain.Protocol) *string {
	if protocol == domain.TCP {
		return aws.String("6")
	}

	if protocol == domain.UDP {
		return aws.String("17")
	}

	if protocol == domain.ICMP {
		return aws.String("1")
	}

	// Fallback to all protocols
	return aws.String("-1")
}

func (a *adapter) buildCreateAclEntryRequest(rule domain.Rule) *ec2.CreateNetworkAclEntryInput {
	input := ec2.CreateNetworkAclEntryInput{
		Egress:       aws.Bool(rule.Direction.IsOutbound()),
		NetworkAclId: aws.String(a.NetworkAclId),
		PortRange:    &ec2.PortRange{From: aws.Int64(int64(rule.Port.BeginPort)), To: aws.Int64(int64(rule.Port.EndPort))},
		Protocol:     a.getProtocolNumber(rule.Protocol),
		RuleAction:   "allow",
		RuleNumber:   a.getNextRuleNumber(),
	}

	if rule.IPNet.IP.To4() == nil {
		input.Ipv6CidrBlock = aws.String(rule.IPNet.String())
	} else {
		input.CidrBlock = aws.String(rule.IPNet.String())
	}

	return &input
}

func (a *adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	for _, rule := range rules {
		input := a.buildCreateAclEntryRequest(rule)
		req := a.client.CreateNetworkAclEntryRequest(input)
		_, _ = req.Send(context.TODO()) // TODO error handling
	}

	return domain.AdapterResult{}
}

func (a *adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	currentRules := a.getPersistedRules()
	for _, rule := range rules {
		persistedRule := a.findRuleInPersistedRules(rule, currentRules)
		if persistedRule == nil {
			// todo log
			continue
		}

		input := ec2.DeleteNetworkAclEntryInput{
			Egress:       aws.Bool(rule.Direction.IsOutbound()),
			NetworkAclId: aws.String(a.NetworkAclId),
			RuleNumber:   persistedRule.RuleNumber,
		}

		req := a.client.DeleteNetworkAclEntryRequest(&input)
		_, _ = req.Send(context.TODO()) // TODO error handling
		delete(a.ruleNumbersIndex, rule.String())
	}

	return domain.AdapterResult{}
}

func (a *adapter) getPersistedRules() (rules []ec2.NetworkAclEntry) {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{
			a.NetworkAclId,
			//	todo filter by gatekeeper tag?
		},
	}

	req := a.client.DescribeNetworkAclsRequest(input)
	resp, err := req.Send(context.Background())
	if err != nil {
		// todo log
		return
	}

	// We can assume we will only have 1 network ACL in here because we filtered for only 1
	rules = resp.NetworkAcls[0].Entries
	return
}

func (a *adapter) findRuleInPersistedRules(rule domain.Rule, persistedRules []ec2.NetworkAclEntry) *ec2.NetworkAclEntry {
	for _, persistedRule := range persistedRules {
		cidrs := []string{*persistedRule.CidrBlock, *persistedRule.Ipv6CidrBlock}
		_, found := Find(cidrs, rule.IPNet.IP.String())
		if found {
			return &persistedRule
		}
	}
	return nil
}
