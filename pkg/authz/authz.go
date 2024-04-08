package authz

import "context"

type AuthzService interface {
	Lookup(ctx context.Context, userIdentity string) ([]Permission, error)
	Enforce(ctx context.Context, userIdentity string, resource string, action string) (bool, error)
}

type Permission struct {
	UserIDOrGroup string
	Resource      string
	Action        string
	Effect        string
}
