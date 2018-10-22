// The entire goal for this package is to simply be an extension of
// the already existing Git HTTP(s) server implementation.
//
// Package ui contains a web ui implementation for GoHub.
package ui

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gohub/bare"
	"gohub/ui/routes"
	"gohub/ui/template"
	"net/http"
)

// Init adds the entire GoHub web UI to the provided base Router.
func Init(base bare.Router, domain string) {
	// Parse UI templates
	template.ParseTmpls()

	// Set domain
	base.Host(domain)

	// Handle assests
	assestDir := viper.GetString("gohub.ui.assests")
	fmt.Println(assestDir)
	base.PathPrefix("/assests/").Handler(http.StripPrefix("/assests/", http.FileServer(http.Dir(assestDir))))

	// TODO: Add Middlewares
	// TODO: Middlewares to add - Auth, CSRF
	// TODO: Optional middlewares - Logging, metrics,

	/***** Base Routes: START ******/
	base.Path("/").Methods("GET").HandlerFunc(routes.GetHome)
	/***** Base Routes: END ******/

	/***** User Routes: START ******/
	usr := base.PathPrefix("/{user}/").Subrouter()

	// Handle user base
	usr.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %s to your GoHub page!", vars["user"])
	})

	// TODO: Handle register and login
	usr.Path("/register").Methods("GET", "POST")

	usr.Path("/login").Methods("GET", "POST")

	/***** User Routes: END ******/

	/***** Repo Routes: START ******/
	repo := usr.PathPrefix("/{repo}/").Subrouter()

	repo.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %s to your repo: %s", vars["user"], vars["repo"])
	})

	repo.PathPrefix("/tree/{branch}")
	/***** Repo Routes: END ******/
}
