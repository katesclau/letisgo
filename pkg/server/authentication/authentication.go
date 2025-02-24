package authentication

import (
	"net/http"
	"net/mail"

	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/models"
	"mnesis.com/pkg/server/authorization"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUser returns the user from the request context
func GetUser(r *http.Request) (*models.User, error) {
	r.ParseForm()

	var loginForm = &LoginForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var user = &models.User{
		Username: loginForm.Username,
		Email: mail.Address{
			Name:    loginForm.Username,
			Address: "someuser@somedomain.com",
		},
		Roles: []authorization.AuthorizationRole{
			authorization.SuperAdmin,
		},
	}

	logrus.WithFields(logrus.Fields{
		"user": user,
	}).Trace("[Authentication] User logged in")

	return user, nil
}
