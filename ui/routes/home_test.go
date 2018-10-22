package routes

import (
	"github.com/spf13/viper"
	"gohub/ui/template"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	viper.Set("gohub.ui.templates", filepath.Join(filepath.Dir(wd), "template"))
	template.ParseTmpls()
	os.Exit(m.Run())
}

func TestGetHome(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)

	GetHome(w, req)

	resp := w.Result()
	io.Copy(os.Stdout, resp.Body)
}
