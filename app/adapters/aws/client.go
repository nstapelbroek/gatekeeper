package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// NewAWSClient is a constructor for an ec2.Client with static credentials
func NewAWSClient(accessKey string, secretKey string, awsRegion string) *ec2.Client {
	credProvider := aws.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Credentials = credProvider
	cfg.Region = awsRegion

	return ec2.New(cfg)
}
