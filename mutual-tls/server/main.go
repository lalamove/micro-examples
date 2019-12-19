package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"demo/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lalamove/micro"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Greeter implements GreeterServer
type Greeter struct {
}

// SayHello implements gRPC endpoint "SayHello"
func (s *Greeter) SayHello(
	ctx context.Context,
	req *proto.HelloRequest,
) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "Hello " + req.Name + ", this is greetings from mutual tls micro server",
	}, nil
}

var _ proto.GreeterServer = (*Greeter)(nil) // make sure it implements the interface

var (
	serverName = "server"
	crt        = "certs/server.crt"
	key        = "certs/server.key"
	ca         = "certs/ca.crt"
)

/******************************************************************************
Mutual tls server with certificate authority
******************************************************************************/
func main() {

	reverseProxyFunc := func(
		ctx context.Context,
		mux *runtime.ServeMux,
		grpcHostAndPort string,
		opts []grpc.DialOption,
	) error {
		return proto.RegisterGreeterHandlerFromEndpoint(ctx, mux, grpcHostAndPort, opts)
	}

	// add swagger definition endpoint
	route := micro.Route{
		Method:  "GET",
		Pattern: micro.PathPattern("hello.swagger.json"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			data, _ := ioutil.ReadFile("proto/hello.swagger.json")
			w.Write(data)
		},
	}

	sf := func() {
		log.Println("Server shutting down")
	}

	// init redoc, enable api docs on http://localhost:28888
	redoc := &micro.RedocOpts{
		Up: true,
	}
	redoc.AddSpec("Greeter", "/hello.swagger.json")

	// load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	// create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Fatal(err)
	}

	// append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	// create the TLS configuration
	serverCreds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	clientCreds := credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	s := micro.NewService(
		micro.Debug(true),
		micro.RouteOpt(route),
		micro.ShutdownFunc(sf),
		micro.Redoc(redoc),
		micro.GRPCServerOption(grpc.Creds(serverCreds)),
		micro.GRPCDialOption(grpc.WithTransportCredentials(clientCreds)),
	)
	proto.RegisterGreeterServer(s.GRPCServer, &Greeter{})

	var httpPort, grpcPort uint16
	httpPort = 28888
	grpcPort = 29999
	if err := s.Start(httpPort, grpcPort, reverseProxyFunc); err != nil {
		log.Fatal(err)
	}
}
