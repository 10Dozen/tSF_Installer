package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tSF_Installer/pkg/application"
)

func GetStatusHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Polling status")

		if app.State == application.APP_STATE_DIFF {
			fmt.Println("State is DIFF -> move to diff page")
			tmpl.ExecuteTemplate(w, "diff", NextUnresolvedDiffInfo(app))
			return
		}

		if app.State == application.APP_STATE_DONE {
			tmpl.ExecuteTemplate(w, "done", NewDoneHandlerData(app))
			return
		}

		// -- Show logs
		logs := prepareLogs(app.Logs, -1)
		tmpl.ExecuteTemplate(w, "progress", logs)

		if app.State == application.APP_STATE_ERROR {
			panic("Error happened")
		}
	}
}
