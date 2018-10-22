package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gohub/bare"
	"gohub/cmd/util"
	"gohub/db"
	"os"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "gohub",
	Short: "GoHub is self-hosted Git service",
	Long: `GoHub is mainly a rewrite of Gogs, another sel-hosted Git service.
The reason for the rewrite is to bring a broad implementation to light.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dbCfg := viper.Sub("database")
		if dbCfg == nil {
			return errors.New("missing database config")
		}

		switch dbCfg.GetString("name") {
		case "postgress", "mysql", "sqlite3":
			db.RegisterDB(&db.SQLDB{})
		case "dgraph":
			db.RegisterDB(&db.DgraphDB{})
		default:
			return errors.New("unknown database name")
		}

		return db.Init(dbCfg)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Retrieve git endpoints from viper
		vEndpoints := viper.Get("git.endpoints")
		if vEndpoints == nil {
			return errors.New("gohub: at minimum one endpoint must be set")
		}
		ends := vEndpoints.([]map[string]interface{})

		// Loop over endpoints and pass their config to the bare package
		endpoints := make(map[bare.Endpoint]bare.Configer)
		for _, e := range ends {
			cfg := viper.New()
			for key, val := range e {
				cfg.Set(key, val)
			}

			endpoint, config := bare.NewEndpoint(cfg)
			endpoints[endpoint] = config
		}

		// Start endpoints and wait for any errors
		return util.StartEndpoints(endpoints)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Provide a custom config file.")
	rootCmd.PersistentFlags().BoolP("api", "a", false, "Expose GraphQL API Endpoint")

	rootCmd.PersistentFlags().Bool("no-git", false, "Disable Git protocol access")
	rootCmd.PersistentFlags().Bool("no-http", false, "Disable HTTP access")
	rootCmd.PersistentFlags().Bool("no-https", false, "Disable HTTPS access")
	rootCmd.PersistentFlags().Bool("no-ssh", false, "Disable SSH access")

	viper.BindPFlag("exposeApi", rootCmd.PersistentFlags().Lookup("api"))
	viper.BindPFlag("git", rootCmd.PersistentFlags().Lookup("no-git"))
	viper.BindPFlag("http", rootCmd.PersistentFlags().Lookup("no-http"))
	viper.BindPFlag("https", rootCmd.PersistentFlags().Lookup("no-https"))
	viper.BindPFlag("ssh", rootCmd.PersistentFlags().Lookup("no-ssh"))
}

func initConfig() {
	viper.SetConfigFile("gohub.yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config file:", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.MergeInConfig(); err != nil {
		fmt.Println("Can't read custom config file:", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
