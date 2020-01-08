package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

func (a *Adapter) createIPPermissions(rules []domain.Rule) []ec2.IpPermission {
	permissions := make([]ec2.IpPermission, len(rules))
	for index, rule := range rules {
		permission := ec2.IpPermission{
			IpProtocol: aws.String(rule.Protocol.String()),
			FromPort:   aws.Int64(int64(rule.Port.BeginPort)),
			ToPort:     aws.Int64(int64(rule.Port.EndPort)),
		}

		if rule.IPNet.IP.To4() == nil {
			permission.Ipv6Ranges = []ec2.Ipv6Range{{CidrIpv6: aws.String(rule.IPNet.String())}}
		} else {
			permission.IpRanges = []ec2.IpRange{{CidrIp: aws.String(rule.IPNet.String())}}
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

	req := a.client.AuthorizeSecurityGroupIngressRequest(&input)
	resp, err := req.Send(context.Background())

	result.Error = err
	if resp != nil {
		result.Output = resp.String()
	}

	return
}

// DeleteRules satisfies the domain.Adapter interface
func (a *Adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	input := ec2.RevokeSecurityGroupIngressInput{
		IpPermissions: a.createIPPermissions(rules),
		GroupId:       aws.String(a.securityGroupID),
	}

	req := a.client.RevokeSecurityGroupIngressRequest(&input)
	resp, err := req.Send(context.Background())

	result.Error = err
	if resp != nil {
		result.Output = resp.String()
	}

	return
}
