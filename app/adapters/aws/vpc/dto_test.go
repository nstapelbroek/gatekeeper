package vpc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nstapelbroek/gatekeeper/domain"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
)

func TestNewACLEntryCollectionWithEmptyData(t *testing.T) {
	collection := NewEntryCollection([]ec2.NetworkAclEntry{})
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.ICMP,
		IPNet:     net.IPNet{},
		Port:      domain.PortRange{},
	}

	assert.NotNil(t, collection.entries)
	assert.Empty(t, collection.entries)
	assert.Nil(t, collection.findByRule(rule))
}

func TestNewACLEntryCollectionWithNil(t *testing.T) {
	collection := NewEntryCollection(nil)
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.ICMP,
		IPNet:     net.IPNet{},
		Port:      domain.PortRange{},
	}

	assert.NotNil(t, collection.entries)
	assert.Empty(t, collection.entries)
	assert.Nil(t, collection.findByRule(rule))
}

func TestACLEntryCollectionCanMapIpv4AclToDomainRule(t *testing.T) {
	cidr := "25.25.25.15/24"
	aclRuleNumber := int64(1337)
	ip, ipNet, _ := net.ParseCIDR(cidr)
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.TCP,
		IPNet:     net.IPNet{IP: ip, Mask: ipNet.Mask},
		Port:      domain.PortRange{BeginPort: 20, EndPort: 22},
	}
	aclEntry := ec2.NetworkAclEntry{
		CidrBlock:  &cidr,
		Egress:     aws.Bool(false),
		PortRange:  &ec2.PortRange{From: aws.Int64(20), To: aws.Int64(22)},
		Protocol:   aws.String(strconv.Itoa(domain.TCP.ProtocolNumber())),
		RuleAction: ec2.RuleActionAllow,
		RuleNumber: &aclRuleNumber,
	}

	collection := NewEntryCollection([]ec2.NetworkAclEntry{aclEntry})

	assert.NotNil(t, collection.entries)
	assert.NotEmpty(t, collection.entries)
	assert.Equal(t, &aclEntry, collection.findByRule(rule))
}

func TestACLEntryCollectionCanMapIpv6AclToDomainRule(t *testing.T) {
	cidr := "2002::1234:abcd:ffff:c0a8:101/64"
	ip, ipNet, _ := net.ParseCIDR(cidr)
	aclRuleNumber := int64(634)

	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.TCP,
		IPNet:     net.IPNet{IP: ip, Mask: ipNet.Mask},
		Port:      domain.PortRange{BeginPort: 20, EndPort: 22},
	}
	aclEntry := ec2.NetworkAclEntry{
		Ipv6CidrBlock: &cidr,
		Egress:        aws.Bool(false),
		PortRange:     &ec2.PortRange{From: aws.Int64(20), To: aws.Int64(22)},
		Protocol:      aws.String(strconv.Itoa(domain.TCP.ProtocolNumber())),
		RuleAction:    ec2.RuleActionAllow,
		RuleNumber:    &aclRuleNumber,
	}

	collection := NewEntryCollection([]ec2.NetworkAclEntry{aclEntry})

	assert.NotNil(t, collection.entries)
	assert.NotEmpty(t, collection.entries)
	assert.Equal(t, &aclEntry, collection.findByRule(rule))
}
