package middlewares

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"github.com/spf13/viper"
	"encoding/base64"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getTestSubject() func(handler http.Handler) http.Handler {
	c := viper.New()
	c.SetDefault("http_auth_username", "user")
	c.SetDefault("http_auth_password", "password")
	return MustAuthenticate(c)
}

func TestMustAuthenticateWithWrongPassword(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth("user", "wrongpassword"))
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Internal handler got called instead of returning an unauthorized response")
	})

	middleWare := getTestSubject()
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestMustAuthenticateWithWrongUsername(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth("wronguser", "password"))
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Internal handler got called instead of returning an unauthorized response")
	})

	middleWare := getTestSubject()
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestMustAuthenticateWithWrongUsernameAndPassword(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth("wronguser", "wrongpassword"))
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Internal handler got called instead of returning an unauthorized response")
	})

	middleWare := getTestSubject()
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestMustAuthenticateWithoutAuthenticationHeader(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Internal handler got called instead of returning an unauthorized response")
	})

	middleWare := getTestSubject()
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestSuccessFullAuthentication(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth("user", "password"))
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	middleWare := getTestSubject()
	rr := httptest.NewRecorder()
	handler := middleWare(testHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
