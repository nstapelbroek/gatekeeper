package vpc

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"strconv"
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

func (a *adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	for _, rule := range rules {
		input := a.buildCreateAclEntryRequest(rule)
		req := a.client.CreateNetworkAclEntryRequest(input)
		_, _ = req.Send(context.TODO()) // TODO error handling
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
			NetworkAclId: aws.String(a.NetworkAclId),
			RuleNumber:   ruleNumber,
		}

		req := a.client.DeleteNetworkAclEntryRequest(&input)
		_, _ = req.Send(context.TODO()) // TODO error handling
	}

	return domain.AdapterResult{}
}

func (a *adapter) getNextRuleNumber() *int64 {
	a.ruleStepsTaken = a.ruleStepsTaken + 1
	return aws.Int64(a.startRuleNumber + int64(a.ruleStepSize*a.ruleStepsTaken))
}

func (a *adapter) buildCreateAclEntryRequest(rule domain.Rule) *ec2.CreateNetworkAclEntryInput {
	input := ec2.CreateNetworkAclEntryInput{
		Egress:       aws.Bool(rule.Direction.IsOutbound()),
		NetworkAclId: aws.String(a.NetworkAclId),
		PortRange:    &ec2.PortRange{From: aws.Int64(int64(rule.Port.BeginPort)), To: aws.Int64(int64(rule.Port.EndPort))},
		Protocol:     aws.String(strconv.Itoa(rule.Protocol.ProtocolNumber())),
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

func (a *adapter) getPersistedAclEntries() *aclEntryCollection {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{
			a.NetworkAclId,
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
