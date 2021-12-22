package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradback/grpc-spire/api/adder/v1"
	"github.com/gorilla/mux"
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
	addr = flag.String("addr", "add-service:50051", "the address of adder service")
)

func main() {
	flag.Parse()

	routes := mux.NewRouter()
	routes.HandleFunc("/add/{a}/{b}", handler).Methods("GET")

	fmt.Println("Application is running on: 8080 ...")
	http.ListenAndServe(":8080", routes)
}

func CallAdder(ctx context.Context, in *adder.AddRequest) (*adder.AddResponse, error) {
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketpath)))
	if err != nil {
		log.Fatalf("Unable to create X509Source: %v", err)
	}
	defer source.Close()

	svid, err := source.GetX509SVID()
	if err != nil {
		log.Fatalf("Unable to get X509SVID: %v", err)
	}
	log.Printf("Client SVID: %v", svid.ID)

	// Allowed SPIFFE ID
	serverID := spiffeid.Must("example.org", "ns/default/add")

	conn, err := grpc.DialContext(ctx, *addr, grpc.WithTransportCredentials(
		grpccredentials.MTLSClientCredentials(source, source, tlsconfig.AuthorizeID(serverID)),
	))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return adder.NewAddServiceClient(conn).Compute(ctx, in)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	a, err := strconv.ParseUint(vars["a"], 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid parameter A")
	}
	b, err := strconv.ParseUint(vars["b"], 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid parameter B")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	req := &adder.AddRequest{A: a, B: b}
	if resp, err := CallAdder(ctx, req); err == nil {
		msg := fmt.Sprintf("Summation is %d", resp.Result)
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := fmt.Sprintf("Internal server error: %v", err)
		json.NewEncoder(w).Encode(msg)
	}
}