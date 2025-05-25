package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"time"

	"tSF_Installer/pkg/application"
	"tSF_Installer/pkg/handlers"

	"github.com/gorilla/mux"
)

const (
	DZN_COMMON_FUNCTIONS = "dzn_commonFunctions"
	DZN_GEAR             = "dzn_gear"
	DZN_DYNAI            = "dzn_dynai"
	DZN_TSFRAMEWORK      = "dzn_tSFramework"
	DZN_BRV              = "dzn_brv"

	HTTPS_PREFIX           = "https://"
	GITHUB_DOWNLOAD_SUFFIX = "archive/refs/heads/%s.zip"
	GITHUB_DEFAULT_BRANCH  = "master"

	TEMP_DIR = "FrameworkLatest"
)

var (
	defaults = map[string]string{
		DZN_COMMON_FUNCTIONS: "github.com/10Dozen/dzn_commonFunctions/tree/v1.8", //"github.com/10Dozen/dzn_commonFunctions",
		DZN_GEAR:             "github.com/10Dozen/dzn_gear/tree/v2.11",           //"github.com/10Dozen/dzn_gear",
		DZN_DYNAI:            "github.com/10Dozen/dzn_dynai/tree/v1.3.3",         //"github.com/10Dozen/dzn_dynai",
		DZN_TSFRAMEWORK:      "github.com/10Dozen/dzn_tSFramework/tree/v2.14",    // "github.com/10Dozen/dzn_tSFramework",
	}
	tmpl *template.Template
	app  *application.Application
)

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*html"))
	app = application.NewApplication()
}

func main() {
	go func() {
		time.Sleep(1 * time.Second)
		exec.Command("cmd", "/C", "start", "http:\\localhost:3000").Run()
	}()

	gRouter := mux.NewRouter()
	gRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// Endpoints
	gRouter.HandleFunc("/", handlers.GetRootPageHandler(app, tmpl)).Methods("GET")
	gRouter.HandleFunc("/validate", handlers.GetValidationHandler(app, tmpl)).Methods("POST")
	gRouter.HandleFunc("/install", handlers.GetInstallHandler(app, tmpl)).Methods("POST")
	gRouter.HandleFunc("/status", handlers.GetStatusHandler(app, tmpl)).Methods("GET")
	gRouter.HandleFunc("/diff/{id}", handlers.GetDiffHandler(app, tmpl)).Methods("GET")
	gRouter.HandleFunc("/diff/{id}/{resolution}", handlers.GetDiffResolveHandler(app, tmpl)).Methods("PATCH")
	gRouter.HandleFunc("/diff/skipAll", handlers.GetDiffSkipAllHandler(app, tmpl)).Methods("POST")

	fmt.Println("Сервер запущен по адресу http://localhost:3000/")
	http.ListenAndServe(":3000", gRouter)
}
