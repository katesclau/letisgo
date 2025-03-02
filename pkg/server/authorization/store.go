package authorization

type AuthorizationStore interface {
	Get(key string) (AuthorizationRole, error)
	Set(key string, role AuthorizationRole) error
}

type StaticAuthorizationStore struct {
	Authorizations map[string]AuthorizationRole
}

func New(authorizationRoles map[string]AuthorizationRole) StaticAuthorizationStore {
	return StaticAuthorizationStore{
		Authorizations: authorizationRoles,
	}
}

func (s StaticAuthorizationStore) Get(key string) (AuthorizationRole, error) {
	return s.Authorizations[key], nil
}

func (s StaticAuthorizationStore) Set(key string, role AuthorizationRole) error {
	s.Authorizations[key] = role
	return nil
}
