package session

import (
	"errors"
	"net/http"

	"github.com/go-session/session/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/katesclau/letisgo/internal/models"
	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/sirupsen/logrus"
)

var (
	ErrEmptyToken   = errors.New("Empty Token")
	ErrUserNotFound = errors.New("User Not Found")
	ErrUnauthorized = errors.New("Unauthorized")
)

type SessionManager interface {
	Create(w http.ResponseWriter, r *http.Request, user *models.User) error
	Get(w http.ResponseWriter, r *http.Request) (*models.User, bool)
}

type RedisSessionManagerConfig struct {
	RedisUrl  string
	JWTSecret []byte
}

type RedisSessionManager struct {
	Config    RedisSessionManagerConfig
	AuthStore authorization.AuthorizationStore
}

func New(cfg RedisSessionManagerConfig, authStore authorization.AuthorizationStore) SessionManager {
	return &RedisSessionManager{
		Config:    cfg,
		AuthStore: authStore,
	}
}

func (s *RedisSessionManager) Create(w http.ResponseWriter, r *http.Request, user *models.User) error {
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

func (s *RedisSessionManager) Get(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	var user *models.User
	authorized := true
	logrus.WithFields(logrus.Fields{
		"path": r.URL.Path,
	}).Trace("[Session][GET]")

	// Get required role for request path
	role, err := s.AuthStore.Get(r.URL.Path)
	logrus.WithFields(logrus.Fields{
		"role": role,
		"err":  err,
	}).Trace("[Session][GET] Role")
	if role == authorization.None || err == authorization.ErrRoleNotDefined {
		logrus.WithFields(logrus.Fields{
			"path": r.URL.Path,
		}).Trace("[Session] No role defined for path")
		return user, authorized
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"path":  r.URL.Path,
		}).Error("[Session] Error getting role for path")
		return user, !authorized
	}

	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("[Session] Error starting session store")
	}

	// Get the JWT from the request
	authHeader := r.Header.Get("Authorization")
	// Extract the JWT from the AuthHeader
	jwt := extractJWT(authHeader)
	if jwt == "" {
		logrus.Error("[Session] No JWT")
		http.Error(w, ErrEmptyToken.Error(), http.StatusUnauthorized)
	}

	// Get the user from the session store
	u, ok := store.Get(jwt)
	if !ok || u == nil {
		logrus.Error("[Session] No User")
		http.Error(w, ErrUserNotFound.Error(), http.StatusUnauthorized)
	}

	// Cast the user to the User struct
	user, ok = u.(*models.User)
	if !ok {
		logrus.Error("[Session] Invalid User")
		http.Error(w, ErrUserNotFound.Error(), http.StatusUnauthorized)
	}

	// Check if the user has the required role
	for _, r := range user.Roles {
		authorized = authorized || r >= role
	}

	if !authorized {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Error("[Session] Unauthorized")
		http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
	}

	return user, authorized
}

func (s *RedisSessionManager) generateJWT(user *models.User) (string, error) {
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
