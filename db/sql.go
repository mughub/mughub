package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// TODO: Add support for SQLite3

// SQLDB represents a general SQL database and is built on top of the
// database/sql and database/sql/driver packages from Go's standard library.
//
// Currently supports:
//		* Postgress  (github.com/lib/pq)
//		* MySQL		 (github.com/go-sql-driver/mysql)
//		* SQLite3	 (?)
type SQLDB struct {
	drvr   *sql.DB
	schema graphql.Schema
}

// Init checks basic config values and then calls sql.Open() to
// establish a connection to your SQL database
func (d *SQLDB) Init(schema graphql.Schema, cfg *viper.Viper) (err error) {
	driver := cfg.GetString("name")
	if driver == "" {
		return ErrUnknownDB{"empty database name provided"}
	}

	err = ErrUnknownDB{driver}
	for _, drvr := range sql.Drivers() {
		if drvr == driver {
			err = nil
		}
	}
	if err != nil {
		return
	}

	dsn := cfg.GetString("dsn")
	if dsn == "" {
		return errors.New("db: dsn must be non-empty")
	}

	d.drvr, err = sql.Open(driver, dsn)
	if err != nil {
		return
	}

	// TODO: Add SQL Resolvers to GraphQL API Schema
	return
}

func (d *SQLDB) Do(ctx context.Context, q string, vars map[string]interface{}) *graphql.Result {
	return graphql.Do(graphql.Params{
		Context:        ctx,
		RequestString:  q,
		VariableValues: vars,
	})
}
