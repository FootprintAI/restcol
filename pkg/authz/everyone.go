package authz

import "context"

var (
	_ AuthzService = &AllowEveryOne{}
)

type AllowEveryOne struct{}

func (a *AllowEveryOne) Lookup(ctx context.Context, userIdentity string) ([]Permission, error) {
	return []Permission{
		Permission{
			UserIDOrGroup: "*",
			Resource:      "*",
			Action:        "*",
			Effect:        "allow",
		},
	}, nil
}

func (a *AllowEveryOne) Enforce(ctx context.Context, userIdentity string, resource string, action string) (bool, error) {
	return true, nil
}
