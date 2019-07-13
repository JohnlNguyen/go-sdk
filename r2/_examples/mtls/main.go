package main

import (
	"go-sdk/certutil"
	"go-sdk/graceful"
	"go-sdk/logger"
	"go-sdk/r2"
	"go-sdk/web"
)

func main() {
	log := logger.All()

	// create the ca
	ca, err := certutil.CreateCA()
	if err != nil {
		log.SyncFatalExit(err)
	}

	caKeyPair, err := ca.GenerateKeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	caPool, err := ca.CertPool()
	if err != nil {
		log.SyncFatalExit(err)
	}

	// create the server certs
	server, err := certutil.CreateServer("mtls-example.local", ca, certutil.OptSubjectCommonName("localhost"))
	if err != nil {
		log.SyncFatalExit(err)
	}
	serverKeyPair, err := server.GenerateKeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	client, err := certutil.CreateClient("mtls-client", ca)
	if err != nil {
		log.SyncFatalExit(err)
	}
	clientKeyPair, err := client.GenerateKeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	serverCertManager, err := certutil.NewCertManagerWithKeyPairs(serverKeyPair, []certutil.KeyPair{caKeyPair}, clientKeyPair)
	if err != nil {
		log.SyncFatalExit(err)
	}

	// create a server
	app := web.New().WithLogger(log).WithBindAddr("127.0.0.1:5000")
	app.WithTLSConfig(serverCertManager.TLSConfig)
	go func() {
		if err := graceful.Shutdown(app); err != nil {
			log.SyncFatalExit(err)
		}
	}()
	<-app.NotifyStarted()

	// make some requests ...

	log.SyncInfof("making a secure request")
	if err := r2.New("https://localhost:5000",
		r2.OptTLSRootCAs(caPool),
		r2.OptTLSClientCert([]byte(clientKeyPair.Cert), []byte(clientKeyPair.Key))).Discard(); err != nil {
		log.SyncFatalExit(err)
	} else {
		log.SyncInfof("secure request success")
	}

	log.SyncInfof("making an insecure request")
	if err := r2.New("https://localhost:5000", r2.OptTLSRootCAs(caPool)).Discard(); err != nil {
		log.SyncFatalExit(err)
	}
}
