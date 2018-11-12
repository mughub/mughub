package bare

import (
	"context"
	"github.com/spf13/viper"
	"net"
)

// Set-up default configuration for SSH server
func init() {
	viper.SetDefault("git.ssh.port", "22")
	viper.SetDefault("git.ssh.auth.password", true)
	viper.SetDefault("git.ssh.auth.pubkey", true)
}

type SSHEndpoint struct{}

func newSSHEndpoint() Endpoint { return &SSHEndpoint{} }

func configSSH(cfg *viper.Viper) ConfigFunc {
	return func() net.Listener {
		return getTCPListener(cfg)

	}
}

func (s *SSHEndpoint) Serve(ctx context.Context, ls net.Listener) error {
	return nil
}
