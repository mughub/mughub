package db

import (
	"context"
	"github.com/spf13/viper"
)

// CockroachDB is a database implementation backed by CockroachDB
type CockroachDB struct {}

func (d *CockroachDB) Init(v *viper.Viper) error {
	return nil
}

func (d *CockroachDB) Query(ctx context.Context) error {
	return nil
}

func (d *CockroachDB) Mutate(ctx context.Context) error {
	return nil
}