package routes

import (
	"context"
	"net/http"

	"github.com/katesclau/letisgo/internal/server/authentication"
	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/katesclau/letisgo/internal/server/endpoints"
	"github.com/sirupsen/logrus"
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
	// TODO: Update the component rendering logic
	// component := navigation.Access()
	// component.Render(ctx, w)
}
