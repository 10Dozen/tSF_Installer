package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"tSF_Installer/pkg/application"

	"github.com/gorilla/mux"
)

type DiffInfo struct {
	Idx         int
	Position    string
	File        string
	DiffContent template.HTML
	HasVSC      bool
	BlockUI     bool
}

func GetDiffHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		diff_idx, has_idx := strconv.Atoi(vars["id"])
		fmt.Println("[DiffHandler] diff_idx=", diff_idx)

		if has_idx != nil {
			diff_idx = 0
		}

		nextDiff := NextUnresolvedDiffInfo(app)
		if nextDiff == nil {
			app.FinishInstallation()
			tmpl.ExecuteTemplate(w, "done", NewDoneHandlerData(app))
			return
		}

		tmpl.ExecuteTemplate(w, "diff", nextDiff)
	}
}

func GetDiffResolveHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		diff_idx, _ := strconv.Atoi(vars["id"])
		resolution, _ := vars["resolution"]

		fmt.Println("[DiffResolveHandler] idx=", diff_idx)
		fmt.Println("[DiffResolveHandler] resolution=", resolution)

		if resolution == application.DIFF_RESOLUTION_VSC {
			app.State = application.APP_STATE_DIFF_VSCODE
			go func() {
				app.ResolveDiffViaVSCode(diff_idx)
			}()

			tmpl.ExecuteTemplate(w, "diff", NextUnresolvedDiffInfo(app))
			return
		}

		app.ResolveDiff(diff_idx, resolution)

		nextDiff := NextUnresolvedDiffInfo(app)
		if nextDiff == nil {
			app.FinishInstallation()
			tmpl.ExecuteTemplate(w, "done", NewDoneHandlerData(app))
			return
		}

		tmpl.ExecuteTemplate(w, "diff", nextDiff)
	}
}

func GetDiffSkipAllHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[DiffSkipAllHandler] Skipped!")
		app.FinishInstallation()
		tmpl.ExecuteTemplate(w, "done", NewDoneHandlerData(app))
	}
}
