package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAWSClientWillSetStaticCredentials(t *testing.T) {
	c := NewAWSClient("somekey", "someSecret", "someregion")

	credValues, _ := c.Credentials.Retrieve()

	assert.IsType(t, aws.StaticCredentialsProvider{}, c.Credentials)
	assert.Equal(t, "somekey", credValues.AccessKeyID)
	assert.Equal(t, "someSecret", credValues.SecretAccessKey)
	assert.False(t, credValues.CanExpire)
}
