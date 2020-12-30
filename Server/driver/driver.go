package driver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

var client = newClient()

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	dialOpts := append([]grpc.DialOption{},
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	d, err := grpc.Dial("localhost:8080", dialOpts...)

	handleError(err)

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func runAlter(schema string) {
	err := client.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	handleError(err)
}

func runMutation(object interface{}) {
	txn := client.NewTxn()
	out, err := json.Marshal(object)
	handleError(err)
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: out})
	handleError(err)
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
