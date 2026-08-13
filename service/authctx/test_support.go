package authctx

import (
	"context"

	"github.com/descope/go-sdk/descope"
)

// WithUserToken returns ctx carrying a user token for the given external ID and
// email. It is intended for tests that exercise authenticated handlers or
// services (for example audited writes that resolve the current user from the
// request context). The token context key is unexported, so this helper is the
// only way to inject a user from outside this package.
func WithUserToken(ctx context.Context, externalID, email string) context.Context {
	return context.WithValue(ctx, ctxUserToken, &descope.Token{
		ID:     externalID,
		Claims: map[string]any{"email": email},
	})
}
