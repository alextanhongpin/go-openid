package gsrv

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Generate private key (.key)
// # Key considerations for algorithm "RSA" ≥ 2048-bit
// openssl genrsa -out server.key 2048
//
// # Key considerations for algorithm "ECDSA" ≥ secp384r1
// # List ECDSA the supported curves (openssl ecparam -list_curves)
// openssl ecparam -genkey -name secp384r1 -out server.key
// Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
// Set the FQDN as localhost:8080

// NewTLS returns a server with graceful shutdown.
func NewTLS(port int, r http.Handler, tlsCert, tlsKey string) <-chan struct{} {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		// CipherSuites: []uint16{
		//         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//         tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		//         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		//         tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		// },
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	idle := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, os.Kill)

		<-sigint

		// Receive interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listener, or context timeout.
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idle)
	}()

	log.Printf("listening to port *:%d. press ctrl + c to cancel.", port)
	if err := srv.ListenAndServeTLS(tlsCert, tlsKey); err != http.ErrServerClosed {
		// Error starting or closing listener.
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
	return idle
}
