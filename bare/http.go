// http.go contains a bare bones HTTP Git + API server
package bare

import (
	"context"
	"net"

	"github.com/spf13/viper"
)

func init() {
	// HTTP
	viper.SetDefault("git.http.port", 80)
	viper.SetDefault("git.http.enabled", true)

	// HTTPS
	viper.SetDefault("git.https.port", 443)
	viper.SetDefault("git.https.enabled", true)
}

// HTTPEndpoint represents the Git HTTP(s) protocol endpoint, as well as, an endpoint
// which can extended to include other functionality e.g. a ui ui
type HTTPEndpoint struct {
	r Router
}

func newHttpEndpoint() Endpoint {
	return &HTTPEndpoint{}
}

func configHttp(cfg *viper.Viper) ConfigFunc {
	return func() net.Listener {
		return getTCPListener(cfg)
	}
}

func (s *HTTPEndpoint) Serve(ctx context.Context, ls net.Listener) error {
	return nil
}
