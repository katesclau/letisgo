package authentication

import (
	"net/http"
	"net/mail"

	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/models"
	"mnesis.com/pkg/server/authorization"
)

type (
	// Move forms to controllers and expose only input types
	LoginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RegisterForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	ResetForm struct {
		Email string `json:"email"`
	}
)

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

// CreateUser creates a new user from the request context
func CreateUser(r *http.Request) (*models.User, error) {
	r.ParseForm()

	var loginForm = &RegisterForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	var user = &models.User{
		Username: loginForm.Username,
		Email: mail.Address{
			Name:    loginForm.Username,
			Address: loginForm.Email,
		},
		Roles: []authorization.AuthorizationRole{
			authorization.SuperAdmin,
		},
	}

	logrus.WithFields(logrus.Fields{
		"user": user,
	}).Trace("[Authentication] User created")

	return user, nil
}

// SendResetEmail sends a reset email to the user
func SendResetEmail(r *http.Request) error {
	r.ParseForm()

	var loginForm = &ResetForm{
		Email: r.FormValue("email"),
	}

	// Find user by email

	// notifications.SendResetEmail(loginForm.email)

	logrus.WithFields(logrus.Fields{
		"email": loginForm.Email,
	}).Trace("[Authentication] Reset email sent")

	return nil
}
