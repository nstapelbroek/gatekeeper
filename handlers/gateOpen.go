package handlers

import (
	"net"
	"net/http"
)

func PostOpen(res http.ResponseWriter, req *http.Request) {
	origin := req.Context().Value("origin")
	iets, assertionSuceeded := origin.(net.IP)
	if !assertionSuceeded {
		panic("context origin was somehow not the expected net.IP type")
	}

	res.Write([]byte(iets.To4().String()))
}
