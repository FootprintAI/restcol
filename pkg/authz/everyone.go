package authz

import "context"

var (
	_ AuthzService = &AllowEveryOne{}
)

type AllowEveryOne struct{}

func (a *AllowEveryOne) Lookup(ctx context.Context, userIdentity UserIdentity) ([]Permission, error) {
	return []Permission{
		Permission{
			UserIDOrGroup: "*",
			Resource:      "*",
			Action:        StringAuthzActionRead,
			Effect:        StringAuthzEffectAllow,
		},
	}, nil
}

func (a *AllowEveryOne) Enforce(ctx context.Context, userIdentity UserIdentity, resource Resource, action AuthzAction) (bool, error) {
	return true, nil
}
