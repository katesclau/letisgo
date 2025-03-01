package session

import (
	"errors"
	"net/http"

	"github.com/go-session/session/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/config"
	"mnesis.com/pkg/models"
	"mnesis.com/pkg/server/authorization"
)

var (
	ErrEmptyToken   = errors.New("Empty Token")
	ErrUserNotFound = errors.New("User Not Found")
	ErrUnauthorized = errors.New("Unauthorized")
)

type SessionManager struct {
	Config *config.Config
}

func New(cfg *config.Config) *SessionManager {
	return &SessionManager{
		Config: cfg,
	}
}

func (s *SessionManager) CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("[Session] Error starting session store")
		return err
	}

	// Generate JWT for the user
	jwt, err := s.generateJWT(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("[Session] Error generating JWT")
	}

	store.Set(jwt, user)
	err = store.Save()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("[Session] Error saving session store")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"user": user,
	}).Trace("[Session] User logged in")

	return nil
}

func (s *SessionManager) Authorize(w http.ResponseWriter, r *http.Request, role authorization.AuthorizationRole) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("[Session] Error starting session store")
		return
	}

	// Get the JWT from the request
	authHeader := r.Header.Get("Authorization")
	// Extract the JWT from the AuthHeader
	jwt := extractJWT(authHeader)
	if jwt == "" {
		logrus.Error("[Authorization] No JWT")
		http.Error(w, ErrEmptyToken.Error(), http.StatusUnauthorized)
	}

	// Get the user from the session store
	u, ok := store.Get(jwt)
	if !ok || u == nil {
		logrus.Error("[Authorization] No User")
		http.Error(w, ErrUserNotFound.Error(), http.StatusUnauthorized)
	}

	// Cast the user to the User struct
	user, ok := u.(*models.User)
	if !ok {
		logrus.Error("[Authorization] Invalid User")
		http.Error(w, ErrUserNotFound.Error(), http.StatusUnauthorized)
	}

	// Check if the user has the required role
	var authorized bool
	for _, r := range user.Roles {
		authorized = authorized || r >= role
	}

	if !authorized {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Error("[Authorization] Unauthorized")
		http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
	}

	logrus.WithFields(logrus.Fields{
		"user": user,
	}).Trace("[Authorization] Authorized")
}

func (s *SessionManager) generateJWT(user *models.User) (string, error) {
	// The JWT should contain the user's ID, username, email, and roles
	// The JWT should be signed using a secret key
	// The JWT should have an expiry time
	// The JWT should be returned as a string

	key := s.Config.JWTSecret
	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":      "my-auth-server",
			"sub":      user.ID,
			"username": user.Username,
			"email":    user.Email.Address,
			"roles":    user.Roles,
		})
	token, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func extractJWT(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
