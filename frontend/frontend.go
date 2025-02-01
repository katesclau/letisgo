package frontend

import (
	"html/template"
	"net/http"

	"mnesis.com/pkg/server/endpoints"
)

type Component struct {
	templateFilename string
	endpoints.APIRouteEndpoint
}

func (c Component) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(c.templateFilename))
		tmpl.Execute(w, nil)
	})
}

var Root = Component{
	templateFilename: "frontend/templates/index.html",
}
