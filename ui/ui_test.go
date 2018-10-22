package ui

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	m := mux.NewRouter()

	m.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)

	wd, _ := os.Getwd()
	viper.Set("gohub.ui.templates", filepath.Join(wd, "template"))
	viper.Set("gohub.ui.assests", filepath.Join(wd, "assests"))
	Init(m, "localhost:8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	if err := srv.ListenAndServe(); err != nil {
		t.Error(err)
	}
}
