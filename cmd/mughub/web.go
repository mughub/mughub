package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/mughub/mughub/bare"
	"github.com/mughub/mughub/ui"
	"github.com/mughub/ssr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	var b bytes.Buffer
	goCmd := exec.Command("go", "env", "GOPATH")
	goCmd.Stdout = &b
	err := goCmd.Run()
	if err != nil {
		panic(err)
	}
	goPath := strings.TrimPrefix(b.String(), "|-")
	goPath = strings.TrimSpace(goPath)
	// TODO: Make this better
	ssrPath := filepath.Join(goPath, "pkg", "mod", "github.com", "mughub", "ssr@v1.0.0")

	// UI
	viper.SetDefault("mughub.ui.domain", "localhost")
	viper.SetDefault("mughub.ui.assests", filepath.Join(ssrPath, "assests"))
	viper.SetDefault("mughub.ui.templates", filepath.Join(ssrPath, "templates"))

	ui.RegisterUI(&ssr.UI{})
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start mughub service with Web UI",
	Long: `Provides a Web UI which implements session management, cookies,
and is highly configurable.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verify HTTP(s) is initialized
		if router == nil {
			return errors.New("http(s) must be specified")
		}

		// Initialize User Interface
		uiCfg := viper.Sub("mughub.ui")
		if uiCfg == nil {
			return errors.New("no ssr user interface config details specified")
		}
		ui.Init(router, uiCfg)

		// Start endpoints
		return bare.ListenAndServe(context.Background(), ends...)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
