package bare

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router represents an HTTP(s) router
type Router interface {
	Get(name string) *mux.Route
	GetRoute(name string) *mux.Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	Headers(pairs ...string) *mux.Route
	Host(tpl string) *mux.Route
	Match(req *http.Request, match *mux.RouteMatch) bool
	MatcherFunc(f mux.MatcherFunc) *mux.Route
	Methods(methods ...string) *mux.Route
	NewRoute() *mux.Route
	Path(tpl string) *mux.Route
	PathPrefix(tpl string) *mux.Route
	Queries(pairs ...string) *mux.Route
	Schemes(schemes ...string) *mux.Route
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Use(mwf ...mux.MiddlewareFunc)
}
