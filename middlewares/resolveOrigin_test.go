package middlewares

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getTestSubjectResolveOrigin(resolveType string) func(handler http.Handler) http.Handler {
	// Disable logrus
	logrus.SetOutput(ioutil.Discard)

	c := viper.New()
	c.SetDefault("resolve_type", resolveType)

	return ResolveOrigin(c)
}

func TestResolveFailedReturnsAServerError(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Handler was called while the request should have been terminated")
	})

	middleWare := getTestSubjectResolveOrigin("default")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestResolveOriginWithDefaultResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "123.123.123.123:12345"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contextOrigin := req.Context().Value(OriginContextKey)
		_, assertionSucceeded := contextOrigin.(net.IP)
		if !assertionSucceeded {
			t.Errorf("Response context was not an instance of net.IP")
		}
	})

	middleWare := getTestSubjectResolveOrigin("default")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResolveOriginWithDefaultResolverToLocalAddress(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "127.0.0.1:12345"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Handler was called while the request should have been terminated")
	})

	middleWare := getTestSubjectResolveOrigin("default")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}
