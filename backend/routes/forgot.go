package routes

import (
	"net/http"

	"github.com/katesclau/letisgo/internal/server/authentication"
	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/katesclau/letisgo/internal/server/endpoints"
	"github.com/sirupsen/logrus"
)

var Forgot = endpoints.RouteEndpoint{
	Handler:           forgotHandler,
	AuthorizationRole: authorization.None,
}

var forgotHandler = func(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	err := authentication.SendResetEmail(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Trace("[Forgot] Recover email sent")

	// TODO: Update the component rendering logic
	// component.Render(ctx, w)
}
