package routes

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"mnesis.com/frontend/components/navigation"
	"mnesis.com/pkg/server/authentication"
	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

var Forgot = endpoints.RouteEndpoint{
	Handler:           forgotHandler,
	AuthorizationRole: authorization.None,
}

var forgotHandler = func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := authentication.SendResetEmail(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Trace("[Forgot] Recover email sent")

	component := navigation.Access()
	component.Render(ctx, w)
}
