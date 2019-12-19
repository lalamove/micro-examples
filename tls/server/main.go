package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"demo/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lalamove/micro"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Greeter implements of GreeterServer
type Greeter struct {
}

// SayHello implements gRPC endpoint "SayHello"
func (s *Greeter) SayHello(
	ctx context.Context,
	req *proto.HelloRequest,
) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "Hello " + req.Name + ", this is greetings from tls micro server",
	}, nil
}

var _ proto.GreeterServer = (*Greeter)(nil) // make sure it implements the interface

var (
	serverName = "server"
	crt        = "certs/server.crt"
	key        = "certs/server.key"
)

/******************************************************************************
tls server with server-side encryption that does not expect client
authentication or credentials
*******************************************************************************/
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

	// init redoc, enable api docs on http://localhost:18888
	redoc := &micro.RedocOpts{
		Up: true,
	}
	redoc.AddSpec("Greeter", "/hello.swagger.json")

	// create the TLS credentials
	serverCreds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	clientCreds, err := credentials.NewClientTLSFromFile(crt, serverName)
	if err != nil {
		log.Fatal(err)
	}

	s := micro.NewService(
		micro.Debug(true),
		micro.RouteOpt(route),
		micro.ShutdownFunc(sf),
		micro.Redoc(redoc),
		micro.GRPCServerOption(grpc.Creds(serverCreds)),
		micro.GRPCDialOption(grpc.WithTransportCredentials(clientCreds)),
	)
	proto.RegisterGreeterServer(s.GRPCServer, &Greeter{})

	// run tls server
	var httpPort, grpcPort uint16
	httpPort = 18888
	grpcPort = 19999
	if err := s.Start(httpPort, grpcPort, reverseProxyFunc); err != nil {
		log.Fatal(err)
	}
}
