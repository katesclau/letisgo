package models

import (
	"net/mail"

	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/nyaruka/phonenumbers"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID       ulid.ULID                         `json:"id"`
	Username string                            `json:"username"`
	Email    mail.Address                      `json:"email"`
	Phone    string                            `json:"phone,omitempty"`
	Roles    []authorization.AuthorizationRole `json:"roles"`

	PhoneNumber phonenumbers.PhoneNumber `json:"-,omitempty"`
}
