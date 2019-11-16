package middleware

import (
	"html/template"
	"net/http"
)

func RenderTemplate(tmpl *template.Template, w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
