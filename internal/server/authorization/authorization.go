package authorization

import "errors"

type AuthorizationRole int

const (
	None AuthorizationRole = iota
	Authenticated
	Admin
	SuperAdmin
)

var ErrRoleNotDefined = errors.New("Role not defined")
