package handlers

import (
	"html/template"
	"net/http"
	"tSF_Installer/pkg/application"
)

// Route Handlers
func GetRootPageHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Reset()

		err := tmpl.ExecuteTemplate(w, "root", GetDefaultOptions())
		if err != nil {
			http.Error(w, "Error exeucting template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
