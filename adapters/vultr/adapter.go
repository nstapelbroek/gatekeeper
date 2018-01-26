package vultr

import (
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
	"net/http"
	"fmt"
	"io/ioutil"
	"errors"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"strings"
)

type adapter struct {
	apiKey          string
	firewallGroupID string
}

// NewVultrAdapter will create a new Vultr adapter object, an unexported type.
func NewVultrAdapter(apiKey string, firewallGroupID string) *adapter {
	a := new(adapter)
	a.apiKey = apiKey
	a.firewallGroupID = firewallGroupID

	return a
}

func (adapter *adapter) dissectSingleRule(rule firewall.Rule) (ipType string, subnetSize string, subnet string) {
	ipType = "v6"
	subnetSize = "128"
	subnet = rule.IP.To16().String()

	if rule.IP.To4() != nil {
		ipType = "v4"
		subnetSize = "32"
		subnet = rule.IP.To4().String()
	}

	return
}

func (adapter *adapter) validateRule(rule firewall.Rule) (err error) {
	if !rule.Port.IsSinglePort() {
		return errors.New("unable to process port-ranges in the Vultr adapter right now")
	}

	if rule.Direction.IsOutbound() {
		return errors.New("cannot create or remove outbound rule's in the Vultr Firewall")
	}

	return
}

func (adapter *adapter) CreateRule(rule firewall.Rule) (err error) {
	err = adapter.validateRule(rule)
	if err != nil {
		return err
	}

	return adapter.createInboundRule(rule)
}

func (adapter *adapter) DeleteRule(rule firewall.Rule) (err error) {
	err = adapter.validateRule(rule)
	if err != nil {
		return err
	}

	return adapter.deleteInboundRule(rule)
}

func (adapter *adapter) doRequest(request *http.Request) (statusCode int, responseBody []byte, err error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	statusCode = response.StatusCode
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	logrus.Debugln(fmt.Sprintf("Response code: '%d', body: '%s'", statusCode, responseBody))

	return
}

func (adapter *adapter) createInboundRule(rule firewall.Rule) (err error) {
	ipType, subnetSize, subnet := adapter.dissectSingleRule(rule)
	ruleRequest := NewRuleCreateRequest(adapter.apiKey, adapter.firewallGroupID, ipType, rule.Protocol.String(), subnet, subnetSize, rule.Port.String())
	statusCode, responseBody, err := adapter.doRequest(ruleRequest.request)
	if err != nil {
		return
	}

	if statusCode == http.StatusPreconditionFailed && string(responseBody) == "Unable to add rule: This rule is already defined" {
		// Functionally the request succeeded, trigger a warning due to potential state issues
		logrus.Warnln(fmt.Sprintf("Tried adding rule for %s on port %s but it was already defined", subnet, rule.Port.String()))
		return nil
	}

	if statusCode != http.StatusOK {
		return errors.New("the Vultr api responded with an unexpected HTTP status code")
	}

	return nil
}

func (adapter *adapter) deleteInboundRule(rule firewall.Rule) (err error) {
	var ruleNumber int
	ruleNumber, err = adapter.deterimeRuleNumber(rule)
	if err != nil {
		logrus.Warningln(err.Error())
	}

	deleteRuleRequest := NewRuleDeleteRequest(adapter.apiKey, adapter.firewallGroupID, ruleNumber)
	statusCode, _, err := adapter.doRequest(deleteRuleRequest.request)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		return errors.New("the Vultr api responded with an unexpected HTTP status code")
	}

	return
}

// deterimeRuleNumber Vultr requires a rule-number for deletion, we fetch all the rules to verify remote config state
func (adapter *adapter) deterimeRuleNumber(localRule firewall.Rule) (ruleNumber int, err error) {
	ipType, _, _ := adapter.dissectSingleRule(localRule)
	listRulesRequest := NewRuleListRequest(adapter.apiKey, adapter.firewallGroupID, ipType)
	statusCode, responseBody, err := adapter.doRequest(listRulesRequest.request)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		return ruleNumber, errors.New("the Vultr api responded with an unexpected HTTP status code")
	}

	deserializedResponse := RuleListResponse{}
	err = json.Unmarshal(responseBody, &deserializedResponse)
	if err != nil {
		return
	}

	for _, externalRule := range deserializedResponse {
		if externalRule.Protocol == strings.ToLower(localRule.Protocol.String()) &&
			externalRule.Subnet == localRule.IP.String() &&
			externalRule.Port == localRule.Port.String() {
			ruleNumber = externalRule.RuleNumber
			return
		}
	}

	return ruleNumber, errors.New("failed resolving correct rule_number at your current vultr configuration")
}
