package driver

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

var client = newClient()

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func RunAlter(schema string) {
	err := client.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	handleError(err)
}

func CreateTransaction() *dgo.Txn {
	txn := client.NewTxn()
	return txn
}

func AddMutationToTransaction(txn *dgo.Txn, object interface{}) {
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	out, err := predicateCaseJSON.Marshal(object)
	handleError(err)
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: out})
	handleError(err)
}

func RunMutation(object interface{}) {
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	txn := client.NewTxn()
	out, err := predicateCaseJSON.Marshal(object)
	handleError(err)
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: out, CommitNow: true})
	handleError(err)
}

func RunQuery(query string) []byte {
	txn := client.NewTxn()
	res, err := txn.Query(context.Background(), query)
	handleError(err)
	return res.GetJson()
}

func handleError(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}
