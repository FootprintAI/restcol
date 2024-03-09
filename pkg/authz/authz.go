package authz

type UserIdentity interface {
	Group() string
	UserID() string
}

type Resource interface {
	URL() string
}

type AuthzAction interface {
	action() string
}

type StringAuthzAction string

func (s StringAuthzAction) action() string {
	return string(s)
}

var (
	StringAuthzActionRead      StringAuthzAction = "read"
	StringAuthzActionReadWrite StringAuthzAction = "readwrite"
	StringAuthzActionDelete    StringAuthzAction = "delete"
)

type AuthzEffect interface {
	effect() string
}

type StringAuthzEffect string

func (s StringAuthzEffect) effect() string {
	return string(s)
}

var (
	StringAuthzEffectAllow StringAuthzEffect = "allow"
	StringAuthzEffectDeny  StringAuthzEffect = "deny"
)

type AuthzService interface {
	Lookup(userIdentity UserIdentity) ([]Permission, error)
	Enforce(userIdentity UserIdentity, resource Resource, action AuthzAction) (bool, error)
}

type Permission struct {
	UserIDOrGroup string
	Resource      string
	Action        StringAuthzAction
	Effect        StringAuthzEffect
}
