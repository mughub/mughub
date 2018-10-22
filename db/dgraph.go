package db

import (
	"context"
	"errors"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// DgraphDB is a database implementation backed by Dgraph
type DgraphDB struct {
	c      *dgo.Dgraph
	schema graphql.Schema
}

func (d *DgraphDB) Init(schema graphql.Schema, cfg *viper.Viper) (err error) {
	// Get Dgraph addr
	addr := cfg.GetString("addr")
	if addr == "" {
		return errors.New("db: dgraph addr must be provided")
	}

	// Connect to Dgraph with gRPC
	d.c, err = d.connect(addr)
	if err != nil {
		return
	}

	// TODO: Add GraphQL Resolvers

	// Setup up Dgraph Schema if not already configured
	return d.setup()
}

func (d *DgraphDB) connect(addr string) (*dgo.Dgraph, error) {
	dc, err := grpc.Dial(addr, grpc.WithInsecure())
	return dgo.NewDgraphClient(api.NewDgraphClient(dc)), err
}

func (d *DgraphDB) setup() error {
	return nil
}

func (d *DgraphDB) Do(ctx context.Context, req string, vars map[string]interface{}) *graphql.Result {
	return graphql.Do(graphql.Params{
		Context:        ctx,
		RequestString:  req,
		VariableValues: vars,
	})
}
