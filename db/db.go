// Package db contains manages for GoHub's data
package db

//go:generate mockgen -destination=./mock/mock.go gohub/db Interface

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("database.name", "sqlite")
	viper.SetDefault("database.dsn", "gohub.db")
}

var dbInterface Interface

// RegisterDB registers the chosen database implementation so other packages can
// simply use the package level API for DB calls and not require any sort of "context"
// for data storage.
func RegisterDB(i Interface) { dbInterface = i }

// Interface represents any storage provider as a GraphQL service.
type Interface interface {
	// Init initializes the database implementation.
	// Each implementation should rely on spf13/viper for config.
	Init(schema graphql.Schema, cfg *viper.Viper) error

	// Do executes the provided GraphQL request
	Do(ctx context.Context, req string, vars map[string]interface{}) *graphql.Result
}

type ErrUnknownDB struct {
	name string
}

func (e ErrUnknownDB) Error() string {
	return "db: unsupported database - " + e.name
}

func Init(schema graphql.Schema, cfg *viper.Viper) error {
	return dbInterface.Init(schema, cfg)
}

func Do(ctx context.Context, req string, vars map[string]interface{}) *graphql.Result {
	return dbInterface.Do(ctx, req, vars)
}
