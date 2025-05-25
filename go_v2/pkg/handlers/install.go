package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"slices"
	"tSF_Installer/pkg/application"
)

// TODO: Full validation of URL to early prevent for errors

func GetValidationHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Validation invoked")
		r.ParseForm()
		state := getInstallOptions(r)

		if !state.Verified {
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmpl.ExecuteTemplate(w, "options", state)
			return
		}

		tmpl.ExecuteTemplate(w, "options", state)
	}
}

func GetInstallHandler(app *application.Application, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("On Install - show progress")
		r.ParseForm()
		state := getInstallOptions(r)

		// -- Final validation
		if !state.Verified {
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmpl.ExecuteTemplate(w, "options", state)
			return
		}

		go app.StartInstallation(state.PathValue, state.Components, state.BackupEnabled)
		err := tmpl.ExecuteTemplate(w, "progress", nil)
		if err != nil {
			http.Error(w, "Error exeucting template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// ---
func GetDefaultOptions() *application.InstallOptions {
	return &application.InstallOptions{
		PathValue:     `F:\Workstation\QSP\1\222`,
		PathError:     "",
		Verified:      true,
		BackupEnabled: true,
		Components: []application.InstallOption{
			application.DefaultOptions[application.SLUG_TSFRAMEWORK],
			application.DefaultOptions[application.SLUG_DYNAI],
			application.DefaultOptions[application.SLUG_GEAR],
			application.DefaultOptions[application.SLUG_COMMON_FUNCTION],
		},
	}
}

func getInstallOptions(r *http.Request) *application.InstallOptions {
	state := GetDefaultOptions()
	state.Verified = true

	// -- Path
	pathInput := r.FormValue(application.SLUG_PATH_INPUT)
	pathError := ""
	if !isValidLocalPath(pathInput) {
		pathError = "Указанной директории не существует."
		state.Verified = false
	}
	state.PathValue = pathInput
	state.PathError = pathError

	// -- Components
	for _, slug := range []string{application.SLUG_COMMON_FUNCTION, application.SLUG_GEAR, application.SLUG_DYNAI, application.SLUG_TSFRAMEWORK} {
		checked := r.FormValue(slug+"CB") == "on" || slug == application.SLUG_COMMON_FUNCTION
		uri := r.FormValue(slug + "URLInput")

		if uri == "" {
			uri = application.DefaultOptions[slug].Value
		}

		componentIdx := slices.IndexFunc(state.Components, func(c application.InstallOption) bool { return c.Slug == slug })
		if componentIdx < 0 {
			continue
		}
		state.Components[componentIdx].Checked = checked
		state.Components[componentIdx].Value = uri
		if checked && !isValidURL(uri, true) {
			state.Components[componentIdx].Error = "Неверный формат URL."
			state.Verified = false
		}

		if slug == application.SLUG_DYNAI && checked {
			componentIdx := slices.IndexFunc(state.Components, func(c application.InstallOption) bool { return c.Slug == application.SLUG_GEAR })
			state.Components[componentIdx].Checked = true
			if !isValidURL(state.Components[componentIdx].Value, true) {
				state.Components[componentIdx].Error = "Неверный формат URL."
			}
		}

		if slug == application.SLUG_TSFRAMEWORK && checked {
			for _, slugRelated := range []string{application.SLUG_GEAR, application.SLUG_DYNAI} {
				componentIdx := slices.IndexFunc(state.Components, func(c application.InstallOption) bool { return c.Slug == slugRelated })
				state.Components[componentIdx].Checked = true
				if !isValidURL(state.Components[componentIdx].Value, true) {
					state.Components[componentIdx].Error = "Неверный формат URL."
				}
			}
		}
	}

	// -- Options
	state.BackupEnabled = r.FormValue("makeBackup") == "on"

	return state
}

func isValidLocalPath(path string) bool {
	fmt.Println("(validateLocalPath)", path)
	return application.PathExists(path)
}

func isValidURL(uri string, fast bool) bool {
	fmt.Println("(validateURL)", uri)

	if _, err := url.ParseRequestURI(uri); err != nil {
		fmt.Println("(validateURL) Invalid URI")
		return false
	}

	if fast {
		fmt.Println("(validateURL) Fast check successful")
		return true
	}

	resp, err := http.Head(uri)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	fmt.Println(resp.StatusCode, resp.Status)
	return resp.StatusCode == 200
}
