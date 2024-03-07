package authz

var (
	_ AuthzService = &AllowEveryOne{}
)

type AllowEveryOne struct{}

func (a *AllowEveryOne) Lookup(userIdentity UserIdentity) ([]Permission, error) {
	return []Permission{
		Permission{
			UserIDOrGroup: "*",
			Resource:      "*",
			Action:        StringAuthzActionReadWrite,
			Effect:        StringAuthzEffectAllow,
		},
	}, nil
}

func (a *AllowEveryOne) Enforce(userIdentity UserIdentity, resource Resource, action AuthzAction) (bool, error) {
	return true, nil
}
