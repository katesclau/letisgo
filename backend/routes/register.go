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

var Register = endpoints.RouteEndpoint{
	Handler:           registerHandler,
	AuthorizationRole: authorization.None,
}

var registerHandler = func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := authentication.CreateUser(r)
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
