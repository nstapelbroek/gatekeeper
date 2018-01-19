package vultr

import (
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
)

type Adapter struct {
	ApiKey          string
	FireWallGroupId string
}

func (adapter Adapter) dissectSingleRule(rule firewall.Rule) (ipType string, subnetSize string, subnet string) {
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

func (adapter Adapter) ProcessRule(rule firewall.Rule) (err error) {
	if !rule.Port.IsSinglePort() {
		panic("cannot create port-ranges yet for Vultr Adapter")
		return
	}

	if rule.Direction.IsOutbound() {
		panic("cannot create outbound rule for Vultr Adapter")
		return
	}

	return adapter.addInboundRule(rule)
}

func (adapter Adapter) addInboundRule(rule firewall.Rule) (err error) {
	ipType, subnetSize, subnet := adapter.dissectSingleRule(rule)
	ruleRequest := NewRuleCreateRequest(adapter.ApiKey, adapter.FireWallGroupId, ipType, rule.Protocol.String(), subnet, subnetSize, rule.Port.String())

	response, err := http.DefaultClient.Do(ruleRequest.request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	logrus.Debugln(fmt.Sprintf("Response code: '%d', body: '%s'", response.StatusCode, responseBody))

	// If the rule is already added we can skip further processing, there should be a timer ready for deletion
	if response.StatusCode == 412 && string(responseBody) == "Unable to add rule: This rule is already defined" {
		return nil
	}

	if response.StatusCode != 200 {
		return errors.New("response code is the expected 200 value")
	}

	deserializedResponse := RuleCreateResponse{}
	err = json.Unmarshal(responseBody, &deserializedResponse)
	return
}
