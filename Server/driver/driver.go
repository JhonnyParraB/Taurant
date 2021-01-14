package driver

import (
	"context"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

var client = newClient()

func newClient() *dgo.Dgraph {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	d, err := grpc.Dial(os.Getenv("DGRAPH_DATABASE_URL"), grpc.WithInsecure())
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
