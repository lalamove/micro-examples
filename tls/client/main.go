package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serverName = "server"
	serverCert = "certs/server.crt"
	clientCert = "certs/client.crt"
	clientKey  = "certs/client.key"
)

/******************************************************************************
Client connect to tls grpc server, note that our demo server's name is "server"
******************************************************************************/
func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	creds, err := credentials.NewClientTLSFromFile(serverCert, serverName)
	if err != nil {
		log.Fatal(err)
	}

	// create a connection with the TLS credentials
	conn, err := grpc.Dial("localhost:19999", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	var client = proto.NewGreeterClient(conn)

	req := &proto.HelloRequest{Name: "tls client"}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.GetMessage())
}
