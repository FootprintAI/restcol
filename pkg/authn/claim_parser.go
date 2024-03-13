package authn

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	authnmiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/authn"
)

type AnnonymousClaimParser struct{}

var (
	_ authnmiddleware.ClaimParser = &AnnonymousClaimParser{}
)

func (n *AnnonymousClaimParser) ParseClaim(ctx context.Context, token string) (jwt.Claims, error) {
	return annonymous{}, nil
}

// annonymous implements jwt.Claims interface
type annonymous struct{}

var (
	_ jwt.Claims = annonymous{}
)

func (a annonymous) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), nil
}

func (a annonymous) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now()), nil
}

func (a annonymous) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)), nil
}

func (a annonymous) GetIssuer() (string, error) {
	return "annonymous", nil
}

func (a annonymous) GetSubject() (string, error) {
	return "annonymous", nil

}
func (a annonymous) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings([]string{"nogroup"}), nil
}
