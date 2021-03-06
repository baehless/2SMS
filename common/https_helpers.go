package common

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"
	"strings"
)

// TODO: handle errors (e.g. empty pool)
func NewCertPoolFromDir(dirPath string) (*x509.CertPool, error) {
	newCertPool := x509.NewCertPool()
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Check extension to match .crt
		if strings.Contains(file.Name(), ".crt") {
			cert, err := ioutil.ReadFile(dirPath + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			newCertPool.AppendCertsFromPEM(cert)
		}
	}
	return newCertPool, nil
}

func CreateHttpsClient(serverCertsDir, clientCert, clientPrivKey string) *http.Client {
	// Load server certificates
	serverCertPool, err := NewCertPoolFromDir(serverCertsDir)
	if err != nil {
		log.Fatal(err)
	}

	// Load client cert and key
	cert, err := tls.LoadX509KeyPair(clientCert, clientPrivKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create HTTPS client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      serverCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	return client
}

func CreateHttpsServer(clientCertsDir, listenInterface, listenPort string, handler http.Handler, clientAuth tls.ClientAuthType) *http.Server {
	// Load client certificates
	clientCertPool, err := NewCertPoolFromDir(clientCertsDir)
	if err != nil {
		log.Println("Unable to create client certificate pool: ", err)
		return nil
	}
	// Require scraper certificate verification
	cfg := &tls.Config{
		ClientAuth: clientAuth,
		ClientCAs:  clientCertPool,
	}
	// Run HTTPS server
	srv := &http.Server{
		Addr:      listenInterface + ":" + listenPort,
		Handler:   handler,
		TLSConfig: cfg,
	}

	return srv
}