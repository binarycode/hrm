package server

import (
	"crypto/tls"
	"net/http"

	"github.com/mgutz/logxi/v1"
	"golang.org/x/crypto/acme/autocert"
)

type Config struct {
	Host  string
	HTTPS bool
}

func Start(config Config) {
	if config.HTTPS {
		serveHTTPS(config.Host)
	} else {
		serveHTTP(config.Host, nil)
	}
}

func serveHTTP(host string, handler http.Handler) {
	address := host + ":80"

	log.Info("Starting HTTP server", "address", address)

	if err := http.ListenAndServe(address, handler); err != http.ErrServerClosed {
		log.Fatal("Unable to start HTTP server", "err", err)
	}
}

func serveHTTPS(host string) {
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache("."),
	}

	tlsConfig := &tls.Config{
		GetCertificate: manager.GetCertificate,
	}

	address := host + ":443"

	server := &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
	}

	go serveHTTP(host, manager.HTTPHandler(nil))

	log.Info("Starting HTTPS server", "address", address)

	if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		log.Fatal("Unable to start HTTPS server", "err", err)
	}
}
