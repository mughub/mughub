// Package db defines how μgHub's data is managed.
package db

//go:generate mockgen -package=dbtest -write_package_comment=false -destination=./dbtest/mock.go github.com/mughub/mughub/db Interface

import (
	"context"
	"github.com/spf13/viper"
	"io"
)

var dbInterface Interface

// RegisterDB registers the chosen database implementation so other packages can
// simply use the package level API for DB calls and not require any sort of "context"
// for data storage.
//
func RegisterDB(i Interface) { dbInterface = i }

// Result represents a GraphQL result.
type Result struct {
	// Data represents the returned data from the database.
	Data interface{}

	// Errors contains any errors encountered when executing the request.
	Errors []error
}

// HasErrors returns true if the result encountered any errors.
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// Interface represents any storage provider as a GraphQL service.
type Interface interface {
	// Init initializes the database implementation.
	// Each implementation should rely on spf13/viper for config.
	//
	Init(schema io.Reader, cfg *viper.Viper) error

	// Do executes the provided GraphQL request
	Do(ctx context.Context, req string, vars map[string]interface{}) *Result
}

// ErrUnknownDB represents initializing an unknown database with the db package.
type ErrUnknownDB struct {
	Name string
}

// Error returns the string formatted representation of ErrUnknownDB.
func (e ErrUnknownDB) Error() string {
	return "db: unsupported database - " + e.Name
}

// Init initializes the registered database.
func Init(schema io.Reader, cfg *viper.Viper) error {
	return dbInterface.Init(schema, cfg)
}

// Do executes a GraphQL request on the registered database.
func Do(ctx context.Context, req string, vars map[string]interface{}) *Result {
	return dbInterface.Do(ctx, req, vars)
}
