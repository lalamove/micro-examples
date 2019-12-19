package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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
	ca         = "certs/ca.crt"
)

/***********************************************************************************************
	Client connect to mutual tls grpc server with certificate authority
************************************************************************************************/
func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	// load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(clientCert, clientKey)
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

	// create the TLS credentials for transport
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial("localhost:29999", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	var client = proto.NewGreeterClient(conn)

	req := &proto.HelloRequest{Name: "mutual tls client"}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.GetMessage())
}
