// Copyright 2022 Bradley Beck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/bradback/grpc-spire/api/adder/v1"
	"github.com/spiffe/go-spiffe/v2/spiffegrpc/grpccredentials"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
)

const (
	socketpath = "unix:///run/spire/sockets/agent.sock"
)

var (
	port = flag.Int("port", 50051, "The sever port")
)

type server struct {
	adder.UnimplementedAddServiceServer
}

func (s *server) Compute(ctx context.Context, in *adder.AddRequest) (*adder.AddResponse, error) {
	log.Printf("Received: %v, %v", in.GetA(), in.GetB())
	return &adder.AddResponse{Result: in.GetA() + in.GetB()}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketpath)))
	if err != nil {
		log.Fatalf("Unable to create X509Source: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("Unable to get X509SVID: %v", err)
	}
	log.Printf("Server SVID: %v", svid.ID)

	// Allowed SPIFFE ID
	clientID := spiffeid.Must("example.org", "ns/default/api")

	s := grpc.NewServer(grpc.Creds(
		grpccredentials.MTLSServerCredentials(source, source, tlsconfig.AuthorizeID(clientID)),
	))

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	adder.RegisterAddServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
