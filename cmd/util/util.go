package util

import (
	"context"
	"gohub/bare"
	"golang.org/x/sync/errgroup"
)

// StartEndpoints is a helper function for starting Bare Endpoints
func StartEndpoints(endpoints map[bare.Endpoint]bare.Configer) error {
	eg, ctx := errgroup.WithContext(context.Background())

	for e, c := range endpoints {
		ls := c.Config()
		et := e
		eg.Go(func() error {
			return et.Serve(ctx, ls) // TODO: Implement error handling
		})
	}

	return eg.Wait()
}
