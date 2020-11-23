package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAWSClientWillSetStaticCredentials(t *testing.T) {
	c := NewAWSClient("somekey", "someSecret", "someregion")

	assert.IsType(t, &ec2.Client{}, c)
}
