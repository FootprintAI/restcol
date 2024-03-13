package middleware

import (
	"context"

	authz "github.com/footprintai/restcol/pkg/authz"
)

type AuthzMiddlwareAdaptor struct {
	authzService authz.AuthzService
}

func NewAuthzMiddlwareAdaptor(authzService authz.AuthzService) *AuthzMiddlwareAdaptor {
	return &AuthzMiddlwareAdaptor{
		authzService: authzService,
	}
}

// Enforce implements authmiddleware.Enforce
func (a *AuthzMiddlwareAdaptor) Enforce(ctx context.Context, subject string, object string, action string) (bool, error) {
	return a.authzService.Enforce(
		ctx,
		rawStringUserId(subject),
		rawStringResource(object),
		getAuthzAction(action),
	)
}

type rawStringUserId string

func (r rawStringUserId) UserIDOrGroup() string {
	return string(r)
}

type rawStringResource string

func (r rawStringResource) URL() string {
	return string(r)
}

func getAuthzAction(s string) authz.StringAuthzAction {
	return authz.NewStringAuthzAction(s, authz.StringAuthzActionRead)
}
