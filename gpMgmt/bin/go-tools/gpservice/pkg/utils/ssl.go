package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Credentials interface {
	LoadServerCredentials() (credentials.TransportCredentials, error)
	LoadClientCredentials() (credentials.TransportCredentials, error)
}

type GpCredentials struct {
	CACertPath     string `json:"caCert"`
	ServerCertPath string `json:"serverCert"`
	ServerKeyPath  string `json:"serverKey"`
	TlsEnabled     bool   `json:"tlsEnabled"`
}

func (c GpCredentials) LoadServerCredentials() (credentials.TransportCredentials, error) {
	if c.TlsEnabled {
		serverCert, err := tls.LoadX509KeyPair(c.ServerCertPath, c.ServerKeyPath)
		if err != nil {
			return nil, fmt.Errorf("could not load server credentials: %w", err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAnyClientCert,
		}
		return credentials.NewTLS(config), nil
	}
	return insecure.NewCredentials(), nil
}

func (c GpCredentials) LoadClientCredentials() (credentials.TransportCredentials, error) {
	if c.TlsEnabled {
		caCert, err := os.ReadFile(c.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("error while loading server certificate: %v", err)
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to add server CA's certificate")
		}

		clientCert, err := tls.LoadX509KeyPair(c.ServerCertPath, c.ServerKeyPath)
		if err != nil {
			return nil, fmt.Errorf("error while loading server certificate: %v", err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      certPool,
		}
		return credentials.NewTLS(config), nil
	}
	return insecure.NewCredentials(), nil

}
