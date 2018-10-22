// Package template contains the templates for GoHub's Web UI.
package template

import (
	"github.com/spf13/viper"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// store contains all the template for GoHub's web ui.
var store = template.New("store")

// ParseTmpls parses the GoHub web ui's template dir.
func ParseTmpls() {
	// Get templates dir and set pattern for sub paths
	tmplDir := viper.GetString("gohub.ui.templates")
	tmplsPattern := filepath.Join(tmplDir, "*", "*.tmpl")

	// Get sub paths
	subPaths, err := filepath.Glob(tmplsPattern)
	if err != nil {
		panic(err)
	}

	// Add root home template
	tmplPaths := make([]string, 1, len(subPaths)+1)
	tmplPaths[0] = filepath.Join(tmplDir, "home.tmpl")
	tmplPaths = append(tmplPaths, subPaths...)

	// Parse templates
	for _, tmplPath := range tmplPaths {
		dir, f := filepath.Split(tmplPath)
		base := filepath.Base(dir)
		name := strings.Split(f, ".")[0]

		tmplStr, err := ioutil.ReadFile(tmplPath)
		if err != nil {
			panic(err)
		}
		store.New(filepath.Join(base, name)).Parse(string(tmplStr))
	}
}

// ExecuteTemplate executes the given template name with the provided data.
func ExecuteTemplate(w io.Writer, name string, data map[string]interface{}) error {
	return store.ExecuteTemplate(w, name, data)
}
