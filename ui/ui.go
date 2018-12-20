// Package ui specifies a web ui interface for μgHub.
package ui

import (
	"github.com/mughub/mughub/bare"
	"github.com/spf13/viper"
)

var uiInterface Interface

// RegisterUI registers a UI for use by
func RegisterUI(i Interface) { uiInterface = i }

// Interface represents a web-based UI for μgHub.
type Interface interface {
	Init(base bare.Router, cfg *viper.Viper)
}

// Init adds the entire μgHub web UI to the provided base Router.
func Init(base bare.Router, cfg *viper.Viper) {
	uiInterface.Init(base, cfg)
}
