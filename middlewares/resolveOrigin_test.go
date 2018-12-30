package middlewares

import (
	"bytes"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getTestSubjectResolveOrigin(resolveType string) func(handler http.Handler) http.Handler {
	// Disable logrus
	//logrus.SetOutput(ioutil.Discard)

	c := viper.New()
	c.SetDefault("resolve_type", resolveType)
	c.SetDefault("resolve_header", "X-Forwarded-For")

	return ResolveOrigin(c)
}

const wantedIP = "123.123.123.123"

// Happy paths
func getHappyPathHandlerFunc(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contextOrigin := req.Context().Value(OriginContextKey)
		ip, assertionSucceeded := contextOrigin.(net.IP)
		if !assertionSucceeded {
			t.Errorf("Response context was not an instance of net.IP")
		}

		if ip.String() != wantedIP {
			t.Errorf("resolved IP was malformed got %v wanted %v", ip.String(), wantedIP)
		}
	})
}

func TestResolveOriginWithDefaultResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = wantedIP + ":12345"
	testHandler := getHappyPathHandlerFunc(t)

	middleWare := getTestSubjectResolveOrigin("default")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResolveOriginWithHeaderResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Forwarded-For", wantedIP)
	testHandler := getHappyPathHandlerFunc(t)

	middleWare := getTestSubjectResolveOrigin("headers")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResolveOriginWithCustomHeaderResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Real-IP", wantedIP)
	testHandler := getHappyPathHandlerFunc(t)

	//logrus.SetOutput(ioutil.Discard)
	c := viper.New()
	c.SetDefault("resolve_type", "headers")
	c.SetDefault("resolve_header", "X-Real-IP")

	middleWare := ResolveOrigin(c)
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResolveOriginWithBodyResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(wantedIP)))
	if err != nil {
		t.Fatal(err)
	}

	testHandler := getHappyPathHandlerFunc(t)

	middleWare := getTestSubjectResolveOrigin("body")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// Exceptionals
func TestNoRemoteAddrWhenUsingDefaultResolver(t *testing.T) {
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

func TestInvalidBodyUsingTheBodyResolver(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("invalid ip")))
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Handler was called while the request should have been terminated")
	})

	middleWare := getTestSubjectResolveOrigin("body")
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
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
