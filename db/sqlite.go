package db

import (
	"context"
	"github.com/spf13/viper"
)

// SQLiteDB implements a database interface backed by SQLite
type SQLiteDB struct {}

func (d *SQLiteDB) Init(v *viper.Viper) error {
	return nil
}

func (d *SQLiteDB) Query(ctx context.Context) error {
	return nil
}

func (d *SQLiteDB) Mutate(ctx context.Context) error {
	return nil
}