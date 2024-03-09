package authz

import "context"

type UserIdentity interface {
	UserIDOrGroup() string
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

func NewStringAuthzAction(s string, defaultAction StringAuthzAction) StringAuthzAction {
	switch s {
	case StringAuthzActionRead.action():
		return StringAuthzActionRead
	case StringAuthzActionWrite.action():
		return StringAuthzActionWrite
	case StringAuthzActionDelete.action():
		return StringAuthzActionDelete
	}
	return defaultAction
}

var (
	StringAuthzActionRead   StringAuthzAction = "GET"
	StringAuthzActionWrite  StringAuthzAction = "POST"
	StringAuthzActionDelete StringAuthzAction = "DELETE"
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
	Lookup(ctx context.Context, userIdentity UserIdentity) ([]Permission, error)
	Enforce(ctx context.Context, userIdentity UserIdentity, resource Resource, action AuthzAction) (bool, error)
}

type Permission struct {
	UserIDOrGroup string
	Resource      string
	Action        StringAuthzAction
	Effect        StringAuthzEffect
}
