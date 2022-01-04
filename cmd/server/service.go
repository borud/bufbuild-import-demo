package main

import (
	"io"
	"log"
	"time"

	"go.buf.build/library/go-grpc/borud/gwp"
)

// Service that implements the streaming API
type Service struct{}

// NewService creates a new Service instance
func NewService() *Service {
	return &Service{}
}

// Connect implementation
func (s *Service) Connect(stream gwp.Gateways_ConnectServer) error {
	// Just send same config request over and over as an example
	go func() {
		for serialCount := int32(1); ; serialCount++ {
			err := stream.Send(&gwp.Packet{
				To: &gwp.Address{Addr: &gwp.Address_B32{B32: 1234567}},
				Payload: &gwp.Packet_Config{
					Config: &gwp.Config{
						Serial:  serialCount,
						WantAck: false,
						Config: map[string]*gwp.Value{
							"foo": {Value: &gwp.Value_Int32Val{Int32Val: 4711}},
						},
					},
				},
			})
			if err != nil {
				log.Printf("unable to send packet: %v", err)
				return
			}

			time.Sleep(2 * time.Second)
		}
	}()

	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("incoming packet: %+v", packet)
	}
}
