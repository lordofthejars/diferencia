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
	Entries []exporter.Entry
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	element := ExtractFile(*r.URL)

	err := renderHtmlTemplate(element, w, DashboardVO{exporter.Entries()}, site)

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
