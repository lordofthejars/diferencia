package core

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/lordofthejars/diferencia/exporter"
)

var site = packr.NewBox("../site")

type DashboardVO struct {
	Entries       []exporter.Entry
	Configuration DiferenciaConfiguration
}

type FailingEntries struct {
	Endpoint     exporter.URLCall
	ErrorDetails []exporter.ErrorData
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	element := ExtractFile(*r.URL)

	err := renderHtmlTemplate(element, w, DashboardVO{exporter.Entries(), *Config}, site)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

}

func dashboardDetailsHandler(w http.ResponseWriter, r *http.Request) {

	method := r.URL.Query().Get("method")
	path := r.URL.Query().Get("path")

	entry := exporter.FindEntry(method, path)
	err := renderHtmlTemplate("diff.html", w, FailingEntries{Endpoint: entry.Endpoint, ErrorDetails: entry.ErrorDetails}, site)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
}

func renderHtmlTemplate(tmplName string, w http.ResponseWriter, p interface{}, box packr.Box) error {

	html, err := box.MustString(tmplName)

	if err != nil {
		return err
	}

	templates, _ := template.New(tmplName).Parse(html)
	return templates.Execute(w, p)

}
