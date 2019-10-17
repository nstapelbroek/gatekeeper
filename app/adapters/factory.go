// Package adapters holds the generic structs and interfaces for the adapter implementations and resolvers
package adapters

import (
	"errors"
	"github.com/nstapelbroek/gatekeeper/app/adapters/aws"
	"github.com/nstapelbroek/gatekeeper/app/adapters/aws/ec2"
	"github.com/nstapelbroek/gatekeeper/app/adapters/digitalocean"
	"github.com/nstapelbroek/gatekeeper/app/adapters/vultr"
	"github.com/nstapelbroek/gatekeeper/domain"
	"github.com/spf13/viper"
)

type AdapterFactory struct {
	config            *viper.Viper
	adapterCollection []domain.Adapter
}

func NewAdapterFactory(config *viper.Viper) (*AdapterFactory, error) {
	f := new(AdapterFactory)
	f.config = config

	doToken := config.GetString("digitalocean_personal_access_token")
	doFirewallId := config.GetString("digitalocean_firewall_id")
	if len(doToken) > 0 && len(doFirewallId) > 0 {
		f.adapterCollection = append(f.adapterCollection, digitalocean.NewDigitalOceanAdapter(doToken, doFirewallId))
	}

	vultrToken := config.GetString("vultr_personal_access_token")
	vultrFirewallId := config.GetString("vultr_firewall_id")
	if len(vultrToken) > 0 && len(vultrFirewallId) > 0 {
		f.adapterCollection = append(f.adapterCollection, vultr.NewVultrAdapter(vultrToken, vultrFirewallId))
	}

	awsKey := config.GetString("aws_access_key")
	awsSecret := config.GetString("aws_secret_key")
	awsRegion := config.GetString("aws_region")
	if len(awsKey) > 0 && len(awsSecret) > 0 && len(awsRegion) > 0 {
		awsClient := aws.NewAWSClient(awsKey, awsSecret, awsRegion)
		if awsSecurityGroupId := config.GetString("aws_security_group_id"); len(awsSecurityGroupId) > 0 {
			f.adapterCollection = append(f.adapterCollection, ec2.NewAWSSecurityGroupAdapter(awsClient, awsSecurityGroupId))
		}
	}

	if len(f.adapterCollection) == 0 {
		return f, errors.New("could not configure any adapters, please set your environment variables")
	}

	return f, nil
}

func (c AdapterFactory) GetAdapters() (adapterCollection *[]domain.Adapter) {
	return &c.adapterCollection
}
