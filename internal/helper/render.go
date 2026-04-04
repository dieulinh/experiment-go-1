package helper

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, data any, files ...string) {
	// always include layout + the specific page files
	templates := append([]string{"templates/layout.html"}, files...)

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "render error: "+err.Error(), 500)
		return
	}
}
