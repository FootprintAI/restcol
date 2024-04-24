package authz

import "context"

var (
	_ AuthzService = &AllowEveryOne{}
)

type AllowEveryOne struct{}

func (a *AllowEveryOne) Enforce(ctx context.Context, userIdentity string, resource string, action string) (bool, error) {
	return true, nil
}
