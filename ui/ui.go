// Package ui specifies a web ui interface for GoHub.
package ui

import (
	"github.com/mughub/mughub/bare"
	"github.com/spf13/viper"
)

var uiInterface Interface

// RegisterUI registers a UI for use by
func RegisterUI(i Interface) { uiInterface = i }

// Interface represents a web-based UI for GoHub.
type Interface interface {
	Init(base bare.Router, cfg *viper.Viper)
}

// Init adds the entire GoHub web UI to the provided base Router.
func Init(base bare.Router, cfg *viper.Viper) {
	uiInterface.Init(base, cfg)
}
