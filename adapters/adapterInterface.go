package adapters

import (
	"net/http"
	"github.com/nstapelbroek/gatekeeper/domain/firewall"
)

type Adapter interface {
	open(rule firewall.Rule) http.Response
	close(rule firewall.Rule) http.Response
}
