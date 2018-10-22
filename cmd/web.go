package cmd

import (
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start GoHub service with Web UI",
	Long: `Provides a Web UI which implements session management, cookies,
and is highly configurable.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement starting GoHub along with Web UI
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
