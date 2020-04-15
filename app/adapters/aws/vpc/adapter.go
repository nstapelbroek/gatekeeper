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
	highestRN := currentEntries.highestNumber
	if highestRN < a.allowedRuleRange.min {
		highestRN = a.allowedRuleRange.min
	}

	if a.allowedRuleRange.max < (highestRN + int64(len(rules))) {
		return domain.AdapterResult{Error: errors.New("not enough rule numbers available")}
	}

	for i, rule := range rules {
		if currentEntries.findByRule(rule) != nil {
			continue
		}

		n := highestRN + int64(i+1)
		req := a.client.CreateNetworkAclEntryRequest(createAddEntryInput(rule, a.networkACLID, n))

		_, err := req.Send(context.TODO())
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
		persistedRule := currentEntries.findByRule(rule)
		if persistedRule == nil {
			continue
		}

		req := a.client.DeleteNetworkAclEntryRequest(createDeleteEntryInput(persistedRule, a.networkACLID))
		_, _ = req.Send(context.TODO())
	}

	return domain.AdapterResult{}
}

func (a *Adapter) getPersistedACLEntries() *EntryCollection {
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
		return NewEntryCollection(nil) // todo log error
	}

	return NewEntryCollection(resp.NetworkAcls[0].Entries) // Assume only 1 result because we filtered
}
