package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nstapelbroek/gatekeeper/domain"
)

// Adapter is a AWS EC2 Security Groups API implementation of the domain.Adapter interface
type Adapter struct {
	client          *ec2.Client
	securityGroupID string
}

// NewAWSSecurityGroupAdapter is a constructor for Adapter
func NewAWSSecurityGroupAdapter(client *ec2.Client, securityGroupID string) *Adapter {
	return &Adapter{
		client:          client,
		securityGroupID: securityGroupID,
	}
}

// ToString satisfies the domain.Adapter interface
func (a *Adapter) ToString() string {
	return "aws-security-group"
}

func (a *Adapter) createIPPermissions(rules []domain.Rule) []types.IpPermission {
	permissions := make([]types.IpPermission, len(rules))
	for index, rule := range rules {
		permission := types.IpPermission{
			IpProtocol: aws.String(rule.Protocol.String()),
			FromPort:   aws.Int32(int32(rule.Port.BeginPort)),
			ToPort:     aws.Int32(int32(rule.Port.EndPort)),
		}

		if rule.IPNet.IP.To4() == nil {
			permission.Ipv6Ranges = []types.Ipv6Range{{CidrIpv6: aws.String(rule.IPNet.String())}}
		} else {
			permission.IpRanges = []types.IpRange{{CidrIp: aws.String(rule.IPNet.String())}}
		}

		permissions[index] = permission
	}

	return permissions
}

// CreateRules satisfies the domain.Adapter interface
func (a *Adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	input := ec2.AuthorizeSecurityGroupIngressInput{
		IpPermissions: a.createIPPermissions(rules),
		GroupId:       aws.String(a.securityGroupID),
	}

	_, err := a.client.AuthorizeSecurityGroupIngress(context.Background(), &input)
	result.Error = err

	return
}

// DeleteRules satisfies the domain.Adapter interface
func (a *Adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	input := ec2.RevokeSecurityGroupIngressInput{
		IpPermissions: a.createIPPermissions(rules),
		GroupId:       aws.String(a.securityGroupID),
	}

	_, err := a.client.RevokeSecurityGroupIngress(context.Background(), &input)
	result.Error = err

	return
}
