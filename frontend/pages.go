package frontend

import (
	"html/template"
	"net/http"

	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

var Pages = endpoints.APIRouteEndpoint{
	Handler:           pagesHandler,
	AuthorizationRole: authorization.None,
}

var pagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/index.html"))
	tmpl.Execute(w, nil)
})
