package bare

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
	"net"
	"os"
)

func init() {
	viper.SetDefault("git.git.port", 9418)
}

type GitEndpoint struct {
	l server.Loader
}

func newGitEndpoint() Endpoint {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return &GitEndpoint{
		l: server.NewFilesystemLoader(osfs.New(wd)),
	}
}

func configGit(cfg *viper.Viper) ConfigFunc {
	return func() net.Listener {
		return getTCPListener(cfg)
	}
}

func (s *GitEndpoint) Serve(ctx context.Context, ls net.Listener) error {
	conn, err := ls.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	refs := packp.NewReferenceUpdateRequest()
	err = refs.Decode(conn)
	fmt.Println(refs)
	return err
}
