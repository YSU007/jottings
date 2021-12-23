// Package main implements a client for Greeter HelloGameX.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"rpcimpl/rpcserver/rpcproto"
)

const (
	address     = "localhost:50051"
	defaultName = "RPC Cli"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpcproto.NewGameXClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r1, err := c.HelloGameX(ctx, &rpcproto.HelloReq{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r1.GetMessage())
	r2, err := c.ByeGameX(ctx, &rpcproto.ByeReq{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r2.GetMessage())
}
