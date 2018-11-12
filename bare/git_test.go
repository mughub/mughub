package bare

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/server"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"os"
	"testing"
)

func TestGitEndpoint_Serve(t *testing.T) {
	cfg := viper.New()
	cfg.Set("port", 9418)

	fs := memfs.New()
	ms := memory.NewStorage()
	repo, err := git.Init(ms, fs)
	if err != nil {
		t.Error(err)
	}
	testF, _ := fs.Create("hello.world")
	testF.Write([]byte("test file for git repo"))
	testF.Close()
	repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"git://localhost/test.git"},
	})
	fmt.Println(repo.CommitObjects())

	l := server.NewFilesystemLoader(fs)
	e, c := &GitEndpoint{l: l}, configGit(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		er := e.Serve(ctx, c())
		if er != nil {
			cancel()
		}
	}()

	fmt.Println(repo.Push(&git.PushOptions{Progress: os.Stdout}))
	fmt.Println("Pushed")

	<-ctx.Done()
}
