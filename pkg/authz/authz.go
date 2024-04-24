package authz

import "context"

type AuthzService interface {
	Enforce(ctx context.Context, userIdentity string, resource string, action string) (bool, error)
}
