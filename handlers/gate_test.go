package handlers

import (
	"context"
	"github.com/nstapelbroek/gatekeeper/adapters/dummy"
	"github.com/nstapelbroek/gatekeeper/middlewares"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func prepareRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, middlewares.OriginContextKey, net.ParseIP("127.0.0.1"))
	req = req.WithContext(ctx)

	return req
}

func TestGateOpenHandler(t *testing.T) {
	req := prepareRequest(t)
	adapterInstance := dummy.NewDummyAdapter()
	gateHandler := NewGateHandler(adapterInstance, 1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gateHandler.PostOpen)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `127.0.0.1 has been whitelisted for 120 seconds`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
