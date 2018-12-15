package bare

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
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

// Endpoint represents a Git protocol serving endpoint.
type Endpoint interface {
	// ListenAndServe is responsible for serving a Git protocol.
	ListenAndServe(ctx context.Context) error
}

// ListenAndServe starts each provided endpoint concurrently.
// If one of them returns an Error it will shutdown the others.
func ListenAndServe(ctx context.Context, ends ...Endpoint) error {
	errChan := make(chan error, len(ends))
	eCtx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	for _, end := range ends {
		go func(e Endpoint) {
			wg.Add(1)
			defer wg.Done()

			err := e.ListenAndServe(eCtx)
			if err != nil {
				cancel()
				errChan <- err
			}
		}(end)
	}

	wg.Wait()
	cancel()

	err := <-errChan
	close(errChan)
	for range errChan {
	} // Drain any other errors from channel

	return err
}
