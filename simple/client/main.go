package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/proto"

	"google.golang.org/grpc"
)

/***********************************************************************************************
	Client connect to insecure grpc server
***********************************************************************************************/
func main() {
	conn, _ := grpc.Dial("localhost:9999", grpc.WithInsecure())
	var client = proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	var req = &proto.HelloRequest{Name: "insecure client"}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.GetMessage())
}
