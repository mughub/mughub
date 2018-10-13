package db

import (
	"context"
	"github.com/spf13/viper"
)

// DgraphDB is a database implementation backed by Dgraph
type DgraphDB struct {}

func (d *DgraphDB) Init(v *viper.Viper) error {
	return nil
}

func (d *DgraphDB) Query(ctx context.Context) error {
	return nil
}

func (d *DgraphDB) Mutate(ctx context.Context) error {
	return nil
}