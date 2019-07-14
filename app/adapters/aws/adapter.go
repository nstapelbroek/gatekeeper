package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
)

type adapter struct {
	client          *ec2.Client
	securityGroupId string
}

// todo: how are you going to handle the different regions?
func NewAWSAdapter(accessKey string, secretKey string, securityGroupId string, awsRegion string) *adapter {
	adapter := new(adapter)
	adapter.securityGroupId = securityGroupId

	credProvider := aws.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Credentials = credProvider
	cfg.Region = awsRegion

	adapter.client = ec2.New(cfg)

	return adapter
}

func (a *adapter) ToString() string {
	return "aws"
}

func (a *adapter) createIpPermissions(rules []domain.Rule) []ec2.IpPermission {
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

func (a *adapter) CreateRules(rules []domain.Rule) (result domain.AdapterResult) {
	input := ec2.AuthorizeSecurityGroupIngressInput{
		IpPermissions: a.createIpPermissions(rules),
		GroupId:       aws.String(a.securityGroupId),
	}

	req := a.client.AuthorizeSecurityGroupIngressRequest(&input)
	resp, err := req.Send(context.Background())

	result.Error = err
	if resp != nil {
		result.Output = resp.String()
	}

	return
}

func (a *adapter) DeleteRules(rules []domain.Rule) (result domain.AdapterResult) {
	input := ec2.RevokeSecurityGroupIngressInput{
		IpPermissions: a.createIpPermissions(rules),
		GroupId:       aws.String(a.securityGroupId),
	}

	req := a.client.RevokeSecurityGroupIngressRequest(&input)
	resp, err := req.Send(context.Background())

	result.Error = err
	if resp != nil {
		result.Output = resp.String()
	}

	return
}
