package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/pedromspeixoto/go-grpc-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "investments-service.staging-15-services.backend.staging.internal.getground.co.uk:50051"
)

func main() {
	/*
		// Load client's certificate and private key
		cert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
		if err != nil {
			log.Fatalf("failed to load client key pair: %s", err)
		}

		// Load the CA certificate
		caCert, err := ioutil.ReadFile("../certs/ca.crt")
		if err != nil {
			log.Fatalf("failed to read CA certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatalf("failed to append CA certificate to pool")
		}

		// Create the TLS credentials for the client
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		})
	*/

	// Load the CA certificate
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %s", err)
	}

	//rootCAs, err := loadSystemRootCAs()
	//if err != nil {
	//	log.Fatalf("failed to load system root CAs: %v", err)
	//}

	// Create TLS credentials with the root CA pool
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      rootCAs,
	})

	// Create a connection to the server
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewYourServiceClient(conn)

	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := c.YourMethod(ctx, &pb.YourRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func loadSystemRootCAs() (*x509.CertPool, error) {
	certFile := "certs/client.crt"

	// Load system root CA certificates
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system root CAs: %v", err)
	}

	// If the system root CA pool is empty, create a new pool
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read custom CA certificates from the specified file
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificates file: %v", err)
	}

	// Append custom CA certificates to the root CA pool
	if ok := rootCAs.AppendCertsFromPEM(certData); !ok {
		return nil, fmt.Errorf("failed to append custom CA certificates to root pool")
	}

	return rootCAs, nil
}
