package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/52north/admission-webhook-server/pkg/admission/admit"
	"github.com/52north/admission-webhook-server/pkg/admission/podnodesselector"
	"github.com/52north/admission-webhook-server/pkg/admission/podtolerationrestriction"
	"github.com/52north/admission-webhook-server/pkg/utils"
)

// TLS secrets
const (
	tlsDir  = `/run/secrets/tls`
	tlsCert = `tls.crt`
	tlsKey  = `tls.key`
)

// Port to listen to
const (
	ENV_LISTEN_PORT = "LISTEN_PORT"
	listenPort      = ":8443"
)

func main() {
	cert := filepath.Join(tlsDir, tlsCert)
	key := filepath.Join(tlsDir, tlsKey)

	mux := http.NewServeMux()
	ctrl := admit.New()
	mux.Handle(admit.GetBasePath(), ctrl)
	log.Print("Registering handlers...")
	registerAllHandlers(ctrl)

	// Config server
	server := &http.Server{
		Addr:    utils.GetEnvVal(ENV_LISTEN_PORT, listenPort),
		Handler: mux,
	}

	// Serve
	log.Print("Starting admission webhook server...")
	log.Fatal(server.ListenAndServeTLS(cert, key))
}

// Register all admission handlers
func registerAllHandlers(ctrl admit.AdmissionController) {
	podnodesselector.Register(ctrl)
	podtolerationrestriction.Register(ctrl)
}
