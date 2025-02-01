package authorization

type AuthorizationRole int

const (
	None AuthorizationRole = iota
	Authenticated
	Admin
	SuperAdmin
)
