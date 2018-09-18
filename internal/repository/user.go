package user

// User represents the operations for the user repository.
type User interface {
	Get(id string) (*oidc.IDToken, bool)
	Put(id string, user *oidc.IDToken)
	Delete(id string)
}
