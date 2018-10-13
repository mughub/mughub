// Package db contains manages for GoHub's data
package db

import (
	"context"
	"github.com/spf13/viper"
)

var dbInterface Interface

// RegisterDB registers the chosen database implementation so other packages can
// simply use the package level API for DB calls and not require any sort of "context"
// for data storage.
func RegisterDB(i Interface) { dbInterface = i }

// Interface represents any storage provider as a GraphQL service.
type Interface interface {
	// Init initializes the database implementation.
	// Each implementation should rely on spf13/viper for config.
	Init(v *viper.Viper) error

	// Query
	Query(ctx context.Context) error

	// Mutate
	Mutate(ctx context.Context) error
}

func Init(v *viper.Viper) error { return dbInterface.Init(v) }

func Query(ctx context.Context) error { return dbInterface.Query(ctx) }

func Mutate(ctx context.Context) error { return dbInterface.Mutate(ctx) }