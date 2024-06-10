package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	pb "github.com/pedromspeixoto/go-grpc-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

// server is used to implement your service.
type server struct {
	pb.UnimplementedYourServiceServer
}

func (s *server) YourMethod(ctx context.Context, in *pb.YourRequest) (*pb.YourResponse, error) {
	return &pb.YourResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("failed to load server key pair: %s", err)
	}

	// Load the CA certificate
	/*
		caCert, err := ioutil.ReadFile("ca.crt")
		if err != nil {
			log.Fatalf("failed to read CA certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatalf("failed to append CA certificate to pool")
		}
	*/

	// Create the TLS credentials for the server
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		//ClientCAs:    caCertPool,
	})

	// Create a gRPC server object
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterYourServiceServer(s, &server{})

	// Listen on the port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Server is listening on %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
