package routes

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"mnesis.com/frontend/components/navigation"
	"mnesis.com/pkg/server/authentication"
	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

var Login = endpoints.RouteEndpoint{
	Handler:           loginHandler,
	AuthorizationRole: authorization.None,
}

var loginHandler = func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := authentication.GetUser(r)
	logrus.WithFields(logrus.Fields{
		"user": user,
		"err":  err,
	}).Trace("[Login] User logged in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx = context.WithValue(ctx, "user", user)
	component := navigation.Access()
	component.Render(ctx, w)
}
