package routes

import (
	"fmt"
	"gohub/ui/template"
	"net/http"
)

// DL retrieves the host domain name and the "Accept-Language" header.
// If "Accept-Language" header is missing, then it defaults to "en".
func DL() TmplCtxFunc {
	return func(req *http.Request, m map[string]interface{}) error {
		m["Domain"] = req.Host
		lang := req.Header.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		m["Lang"] = lang
		return nil
	}
}

// GetHome renders the GoHub home page.
func GetHome(w http.ResponseWriter, req *http.Request) {
	// Get template context data from request
	tmplCtx, _ := GetTmplCtx(req, DL())

	// Execute home page template
	err := template.ExecuteTemplate(w, "template/home", tmplCtx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}
