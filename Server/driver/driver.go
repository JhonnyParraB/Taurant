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

func RunAlter(schema string) error {
	err := client.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	if err != nil {
		return err
	}
	return nil
}

func RunMutation(object interface{}) error {
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	txn := client.NewTxn()
	out, err := predicateCaseJSON.Marshal(object)
	if err != nil {
		return err
	}
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: out, CommitNow: true})
	if err != nil {
		return err
	}
	return nil
}

func RunMutationForDelete(object interface{}) error {
	var predicateCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	txn := client.NewTxn()
	out, err := predicateCaseJSON.Marshal(object)
	if err != nil {
		return err
	}
	_, err = txn.Mutate(context.Background(), &api.Mutation{DeleteJson: out, CommitNow: true})
	if err != nil {
		return err
	}
	return nil
}

func RunQuery(query string) ([]byte, error) {
	txn := client.NewTxn()
	res, err := txn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return res.GetJson(), nil
}
