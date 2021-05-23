package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// NewAWSClient is a constructor for an ec2.Client with static credentials
func NewAWSClient(accessKey string, secretKey string, awsRegion string) *ec2.Client {
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
				Source:          "gatekeeper-configuration-arguments",
			},
		}),
	)

	return ec2.NewFromConfig(cfg)
}
