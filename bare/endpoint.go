package bare

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net"
)

type Endpoint interface {
	Serve(ctx context.Context, ls net.Listener) error
}

// NewEndpoint handles all the business logic for setting up and configuring
// Git endpoints.
func NewEndpoint(cfg *viper.Viper) (e Endpoint, c Configer) {
	switch cfg.GetString("protocol") {
	case "git":
		e, c = newGitEndpoint(), configGit(cfg)
	case "ssh":
		e, c = newSSHEndpoint(), configSSH(cfg)
	case "http":
		e, c = newHttpEndpoint(), configHttp(cfg)
	case "https":
		e, c = newHttpEndpoint(), configHttp(cfg)
	}
	return
}

func getTCPListener(cfg *viper.Viper) net.Listener {
	addr := fmt.Sprintf("%s:%d", cfg.GetString("addr"), cfg.GetInt("port"))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		l.Close()
		panic(err)
	}
	return l
}

type Configer interface {
	Config() net.Listener
}

type ConfigFunc func() net.Listener

func (f ConfigFunc) Config() net.Listener {
	return f()
}
