package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mughub/dgraph"
	"github.com/mughub/git"
	"github.com/mughub/gqlite"
	"github.com/mughub/http"
	"github.com/mughub/mughub/api"
	"github.com/mughub/mughub/bare"
	"github.com/mughub/mughub/db"
	"github.com/mughub/sql"
	"github.com/mughub/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	http3 "net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func init() {
	// API
	viper.SetDefault("mughub.endpoint.api.domain", "localhost")

	// Database
	viper.SetDefault("mughub.database.name", "gqlite")

	// Git Protocol
	viper.SetDefault("mughub.endpoint.git.port", 9418)

	// HTTP(s) Protocol
	viper.SetDefault("mughub.endpoint.http.port", 80)             // HTTP port
	viper.SetDefault("mughub.endpoint.http.secure.port", 433)     // HTTPS port
	viper.SetDefault("mughub.endpoint.http.redirect", false)      // Redirect HTTP to HTTPS

	// SSH Protocol
	viper.SetDefault("mughub.endpoint.ssh.port", 22)
	viper.SetDefault("mughub.endpoint.ssh.auth.password", true)
	viper.SetDefault("mughub.endpoint.ssh.auth.pubkey", true)
}

var (
	cfgFile string
	ends    []bare.Endpoint
	router  bare.Router
)

// setupDB handles identifying the desired DB implementation and register/initializing
func setupDB(cfg *viper.Viper) error {
	switch cfg.GetString("name") {
	case "postgress", "mysql":
		db.RegisterDB(&sql.DB{})
	case "dgraph":
		db.RegisterDB(&dgraph.DB{})
	case "gqlite":
		db.RegisterDB(&gqlite.DB{})
	default:
		return errors.New("unknown database name")
	}

	return db.Init(api.Schema, cfg)
}

// getEnds identifies and creates all desired endpoints
func getEnds() bare.Router {
	if viper.GetBool("mughub.endpoint.git.enabled") {
		e := git.NewEndpoint()
		ends = append(ends, e)
	}

	if viper.GetBool("mughub.endpoint.ssh.enabled") {
		e := ssh.NewEndpoint()
		ends = append(ends, e)
	}

	if viper.GetBool("mughub.endpoint.http.enabled") {
		e, r := http.NewEndpoint(viper.Sub("mughub.endpoint.http"))
		ends = append(ends, e)
		return r
	}

	return nil
}

type apiEndpoint struct {
	s *http3.Server
}

func (e *apiEndpoint) ListenAndServe(ctx context.Context) error {
	return e.s.ListenAndServe()
}

var rootCmd = &cobra.Command{
	Use:   "mughub",
	Short: "μghub is self-hosted Git service",
	Long: `μghub is designed to provide a highly flexible Git service. This
root command will launch GoHub as a bare bones service with Git protocol
endpoints, a database and an API endpoint. It will NOT launch with a UI. In
order to launch with a UI, see the web sub command.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" {
			return nil
		}

		dbCfg := viper.Sub("mughub.database")
		if dbCfg == nil {
			return errors.New("missing database config")
		}

		err := setupDB(dbCfg)
		if err != nil {
			return err
		}

		endCfg := viper.Sub("mughub.endpoint")
		if endCfg == nil {
			return errors.New("missing git endpoints config")
		}

		router = getEnds()
		apiCfg := endCfg.Sub("api")
		if apiCfg != nil {
			if router == nil {
				router = mux.NewRouter()
				apiEnd := &apiEndpoint{
					s: &http3.Server{
						Addr: ":8080",
						Handler: router,
						// TODO: Set timeouts and only serve over HTTPS
					},
				}
				ends = append(ends, apiEnd)
			}

			http.RegisterAPIEndpoint(router, apiCfg)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Start endpoints and wait for any errors
		return bare.ListenAndServe(context.Background(), ends...)
	},
}

////// Any stuff below has to do with command configuration and not the actual
////// execution of the root command.

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Provide a custom config file")

	rootCmd.PersistentFlags().Bool("api", false, "Expose GraphQL API Endpoint")
	rootCmd.PersistentFlags().Bool("git", false, "Enable Git protocol access")
	rootCmd.PersistentFlags().Bool("http", false, "Enable HTTP access")
	rootCmd.PersistentFlags().Bool("https", false, "Enable HTTPS access")
	rootCmd.PersistentFlags().Bool("ssh", false, "Enable SSH access")

	err := viper.BindPFlag("mughub.endpoint.api.enabled", rootCmd.PersistentFlags().Lookup("api"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("mughub.endpoint.git.enabled", rootCmd.PersistentFlags().Lookup("git"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("mughub.endpoint.http.enabled", rootCmd.PersistentFlags().Lookup("http"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("mughub.endpoint.http.secure.enabled", rootCmd.PersistentFlags().Lookup("https"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("mughub.endpoint.ssh.enabled", rootCmd.PersistentFlags().Lookup("ssh"))
	if err != nil {
		panic(err)
	}

}

func initConfig() {
	var noCfg bool

	// First, try to read in pre-existing gohub.yml
	viper.SetConfigFile("mughub.yml")
	if err := viper.ReadInConfig(); os.IsNotExist(err) {
		noCfg = true
		fmt.Println("mughub: config file, gohub.yml, doesn't exist already. one will be created later.")
	} else if err != nil {
		fmt.Println("mughub: can't read config file:", err)
	}

	// Next, merge overloaded config vals from a custom config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Println("mughub: can't read custom config file:", err)
			os.Exit(1)
		}
	}

	// Lastly, write gohub.yml if it didn't exist prior
	if noCfg {
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("mughub: couldn't write config to:", cfgFile, err)
			os.Exit(1)
		} else {
			fmt.Println("mughub: created config file, gohub.yml")
		}
	}
}
