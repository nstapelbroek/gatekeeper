package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSClientWillSetStaticCredentials(t *testing.T) {
	c := NewAWSClient("somekey", "someSecret", "someregion")

	assert.IsType(t, &ec2.Client{}, c)
}
