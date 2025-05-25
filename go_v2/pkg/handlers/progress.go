package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"tSF_Installer/pkg/application"
)

var tagToBadge = map[string]string{
	application.LOG_TAG_INFO:  "secondary",
	application.LOG_TAG_ERROR: "danger",
}

type LogData struct {
	Message template.HTML
	Tag     string
	Badge   string
}

func GetStatusHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Polling status")

		if app.State == application.APP_STATE_DIFF {
			fmt.Println("State is DIFF -> move to diff page")
			tmpl.ExecuteTemplate(w, "diff", NextUnresolvedDiffInfo(app))
			return
		}

		if app.State == application.APP_STATE_DONE {
			tmpl.ExecuteTemplate(w, "done", template.HTML(app.TargetDir))
			return
		}

		// -- Show logs
		logs := []LogData{}
		for _, log := range app.Logs {
			logs = append(logs, LogData{template.HTML(log.Message), log.Level, tagToBadge[log.Level]})
		}
		slices.Reverse(logs)

		tmpl.ExecuteTemplate(w, "progress", logs)

		if app.State == application.APP_STATE_ERROR {
			panic("Error happened")
		}
	}
}
