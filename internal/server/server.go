// Package server create http transport
//
// For example create:
//
//	serv := server.NewServer(log, serverAddress, r)
//	serv.Start()
//	serv.Stop()
package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server http server
type Server struct {
	log    *zap.Logger
	errors chan error
	server http.Server
}

// NewServer create http transport
func NewServer(log *zap.Logger, address string, router http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    address,
			Handler: router,
		},
		log:    log,
		errors: make(chan error),
	}
}

// Start http server and listening port
func (s *Server) Start(enabledHTTS bool) {
	go func() {
		s.log.Info("Http Server listening on ", zap.String("addr", s.server.Addr))
		defer close(s.errors)

		if enabledHTTS {
			certPEM, privateKeyPEM, err := s.generateFilesForTLS()

			if err != nil {
				s.errors <- err
				return
			}

			certificate, err := tls.X509KeyPair(certPEM, privateKeyPEM)

			if err != nil {
				s.errors <- err
				return
			}

			tlsConfig := &tls.Config{}
			tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
			s.server.TLSConfig = tlsConfig

			s.errors <- s.server.ListenAndServeTLS("", "")

			return
		}

		s.errors <- s.server.ListenAndServe()
	}()
}

// Stop don't listen new connection and waiting close current active connection. Used for graceful shutdown
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// Notify return chan with error. For example include error "port already exists"
func (s *Server) Notify() <-chan error {
	return s.errors
}

func (s *Server) generateFilesForTLS() ([]byte, []byte, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		s.log.Error("can't create rsa key", zap.Error(err))

		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		s.log.Error("can't create certification", zap.Error(err))

		return nil, nil, err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certPEM.Bytes(), privateKeyPEM.Bytes(), nil
}
