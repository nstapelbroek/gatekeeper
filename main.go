package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"

	"github.com/nstapelbroek/gatekeeper/application"
)

func newConfig() (*viper.Viper, error) {
	c := viper.New()
	c.SetDefault("http_port", "8080")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")
	c.SetDefault("http_auth_username", "user")
	c.SetDefault("http_auth_password", "password")
	c.SetDefault("resolve_type", "RemoteAddr")
	c.SetDefault("resolve_header", "X-Forwarded-For")
	c.SetDefault("rule_close_timeout", 120)
	c.SetDefault("rule_ports", "TCP:22")

	c.AutomaticEnv()

	return c, nil
}

func main() {
	config, err := newConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	app, err := application.New(config)
	if err != nil {
		logrus.Fatal(err)
	}

	logParameter := flag.String("log-level", "info", "determine the verbosity of the logger")
	flag.Parse()
	logLevel, err := logrus.ParseLevel(*logParameter)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetLevel(logLevel)

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	serverAddress := ":" + config.GetString("http_port")
	certFile := config.GetString("http_cert_file")
	keyFile := config.GetString("http_key_file")
	drainIntervalString := config.GetString("http_drain_interval")

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
