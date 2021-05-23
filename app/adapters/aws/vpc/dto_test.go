package vpc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nstapelbroek/gatekeeper/domain"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
)

func TestNewACLEntryCollectionWithEmptyData(t *testing.T) {
	collection := NewEntryCollection([]types.NetworkAclEntry{})
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.ICMP,
		IPNet:     net.IPNet{},
		Port:      domain.PortRange{},
	}

	_, found := collection.findByRule(rule)

	assert.NotNil(t, collection.entries)
	assert.Empty(t, collection.entries)
	assert.False(t, found)
}

func TestNewACLEntryCollectionWithNil(t *testing.T) {
	collection := NewEntryCollection(nil)
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.ICMP,
		IPNet:     net.IPNet{},
		Port:      domain.PortRange{},
	}

	_, found := collection.findByRule(rule)

	assert.NotNil(t, collection.entries)
	assert.Empty(t, collection.entries)
	assert.False(t, found)
}

func TestACLEntryCollectionCanMapIpv4AclToDomainRule(t *testing.T) {
	cidr := "25.25.25.15/24"
	aclRuleNumber := int32(1337)
	ip, ipNet, _ := net.ParseCIDR(cidr)
	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.TCP,
		IPNet:     net.IPNet{IP: ip, Mask: ipNet.Mask},
		Port:      domain.PortRange{BeginPort: 20, EndPort: 22},
	}
	aclEntry := types.NetworkAclEntry{
		CidrBlock:  &cidr,
		Egress:     aws.Bool(false),
		PortRange:  &types.PortRange{From: aws.Int32(20), To: aws.Int32(22)},
		Protocol:   aws.String(strconv.Itoa(domain.TCP.ProtocolNumber())),
		RuleAction: types.RuleActionAllow,
		RuleNumber: &aclRuleNumber,
	}

	collection := NewEntryCollection([]types.NetworkAclEntry{aclEntry})
	entry, found := collection.findByRule(rule)

	assert.NotNil(t, collection.entries)
	assert.NotEmpty(t, collection.entries)
	assert.True(t, found)
	assert.Equal(t, aclEntry, entry)
}

func TestACLEntryCollectionCanMapIpv6AclToDomainRule(t *testing.T) {
	cidr := "2002::1234:abcd:ffff:c0a8:101/64"
	ip, ipNet, _ := net.ParseCIDR(cidr)
	aclRuleNumber := int32(634)

	rule := domain.Rule{
		Direction: domain.Inbound,
		Protocol:  domain.TCP,
		IPNet:     net.IPNet{IP: ip, Mask: ipNet.Mask},
		Port:      domain.PortRange{BeginPort: 20, EndPort: 22},
	}
	aclEntry := types.NetworkAclEntry{
		Ipv6CidrBlock: &cidr,
		Egress:        aws.Bool(false),
		PortRange:     &types.PortRange{From: aws.Int32(20), To: aws.Int32(22)},
		Protocol:      aws.String(strconv.Itoa(domain.TCP.ProtocolNumber())),
		RuleAction:    types.RuleActionAllow,
		RuleNumber:    &aclRuleNumber,
	}

	collection := NewEntryCollection([]types.NetworkAclEntry{aclEntry})
	entry, found := collection.findByRule(rule)

	assert.NotNil(t, collection.entries)
	assert.NotEmpty(t, collection.entries)
	assert.True(t, found)
	assert.Equal(t, aclEntry, entry)
}
