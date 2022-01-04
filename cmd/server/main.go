package main

import (
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
	"go.buf.build/library/go-grpc/borud/gwp"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCAddr string `long:"grpc-addr" default:":5011" description:"gRPC listen address"`
}

func main() {
	_, err := flags.ParseArgs(&opt, os.Args)
	if err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	// Create listen socket
	listenSocket, err := net.Listen("tcp", opt.GRPCAddr)
	if err != nil {
		log.Fatalf("Unable to listen to address %s: %v", opt.GRPCAddr, err)
	}

	// Create the service instance
	service := NewService()

	// Create the gRPC server and start it
	grpcServer := grpc.NewServer()
	gwp.RegisterGatewaysServer(grpcServer, service)
	err = grpcServer.Serve(listenSocket)
	if err != nil {
		log.Fatalf("Serve returned error: %v", err)
	}
}
