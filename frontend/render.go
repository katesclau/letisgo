package frontend

import (
	"net/http"

	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

var Index = endpoints.APIRouteEndpoint{
	Handler:           indexHandler,
	AuthorizationRole: authorization.None,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	component := index("Name")
	component.Render(r.Context(), w)
}
