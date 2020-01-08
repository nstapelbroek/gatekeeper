package adapters

import (
	"errors"
	"github.com/nstapelbroek/gatekeeper/app/adapters/aws"
	"github.com/nstapelbroek/gatekeeper/app/adapters/aws/ec2"
	"github.com/nstapelbroek/gatekeeper/app/adapters/aws/vpc"
	"github.com/nstapelbroek/gatekeeper/app/adapters/digitalocean"
	"github.com/nstapelbroek/gatekeeper/app/adapters/vultr"
	"github.com/nstapelbroek/gatekeeper/domain"
	"github.com/spf13/viper"
)

// AdapterFactory will act as the owner of all adapter instances
type AdapterFactory struct {
	adapterCollection []domain.Adapter
}

// NewAdapterFactory is a constructor method for AdapterFactory
func NewAdapterFactory(config *viper.Viper) (*AdapterFactory, error) {
	f := new(AdapterFactory)

	doToken := config.GetString("digitalocean_personal_access_token")
	doFirewallID := config.GetString("digitalocean_firewall_id")
	if len(doToken) > 0 && len(doFirewallID) > 0 {
		f.adapterCollection = append(f.adapterCollection, digitalocean.NewDigitalOceanAdapter(doToken, doFirewallID))
	}

	vultrToken := config.GetString("vultr_personal_access_token")
	vultrFirewallID := config.GetString("vultr_firewall_id")
	if len(vultrToken) > 0 && len(vultrFirewallID) > 0 {
		f.adapterCollection = append(f.adapterCollection, vultr.NewVultrAdapter(vultrToken, vultrFirewallID))
	}

	awsKey := config.GetString("aws_access_key")
	awsSecret := config.GetString("aws_secret_key")
	awsRegion := config.GetString("aws_region")
	if len(awsKey) > 0 && len(awsSecret) > 0 && len(awsRegion) > 0 {
		awsClient := aws.NewAWSClient(awsKey, awsSecret, awsRegion)
		if awsSecurityGroupID := config.GetString("aws_security_group_id"); len(awsSecurityGroupID) > 0 {
			f.adapterCollection = append(f.adapterCollection, ec2.NewAWSSecurityGroupAdapter(awsClient, awsSecurityGroupID))
		}
		if awsNetworkACLId := config.GetString("aws_network_acl_id"); len(awsNetworkACLId) > 0 {
			f.adapterCollection = append(f.adapterCollection, vpc.NewAWSNetworkACLAdapter(awsClient, awsNetworkACLId, config.GetString("aws_network_acl_rule_number_range")))
		}
	}

	if len(f.adapterCollection) == 0 {
		return f, errors.New("could not configure any adapters, please set your environment variables")
	}

	return f, nil
}

// GetAdapters exposes the internal collection of build adapter instances
func (c AdapterFactory) GetAdapters() (adapterCollection *[]domain.Adapter) {
	return &c.adapterCollection
}
