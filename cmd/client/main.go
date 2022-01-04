package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"go.buf.build/library/go-grpc/borud/gwp"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCAddr string `long:"grpc-addr" default:":5011" description:"gRPC server address"`
}

func main() {
	// Parse command line flags
	_, err := flags.ParseArgs(&opt, os.Args)
	if err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	log.Printf("addr: %s", opt.GRPCAddr)

	// Connect to gRPC server
	conn, err := grpc.Dial(opt.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing gRPC server: %v", err)
	}
	defer conn.Close()

	// Create client
	client := gwp.NewGatewaysClient(conn)

	// Create bi-directional stream
	stream, err := client.Connect(context.Background())
	if err != nil {
		log.Fatalf("error connecting stream: %+v", err)
	}

	// Fire goroutine that listens for incoming messages
	go func() {
		for {
			packet, err := stream.Recv()
			if err == io.EOF {
				log.Printf("server closed connection")
				return
			}
			if err != nil {
				log.Fatalf("failed to receive packet: %v", err)
			}

			log.Printf("> %+v", packet)
		}
	}()

	serialCount := uint32(0)
	for {
		serialCount++

		// Create a packet
		packet := &gwp.Packet{
			Payload: &gwp.Packet_Sample{
				Sample: &gwp.Sample{
					From:      &gwp.Address{Addr: &gwp.Address_B32{B32: 1234567}},
					Timestamp: uint64(time.Now().UnixMilli()),
					Type:      1,
					Value:     &gwp.Value{Value: &gwp.Value_FloatVal{FloatVal: 3.14}},
					Serial:    serialCount,
					WantAck:   false,
				},
			},
		}

		// Send the packet
		err := stream.Send(packet)
		if err != nil {
			log.Fatalf("failed to send packet: %v", err)
		}

		// Wait a bit
		time.Sleep(time.Second)
	}
}
