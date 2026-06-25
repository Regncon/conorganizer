package authctx

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"
)

func TestWithUserToken_StoresTokenRetrievableFromContext(t *testing.T) {
	token := &descope.Token{ID: "external-123", Claims: map[string]any{"email": "player@example.com"}}

	ctx := WithUserToken(context.Background(), token)

	got, err := GetUserTokenFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != token {
		t.Fatalf("token mismatch\nexpected: %v\nactual:   %v", token, got)
	}
}

func TestWithUserToken_ClearsPriorSessionError(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxSessionError, context.DeadlineExceeded)
	token := &descope.Token{ID: "external-123"}

	ctx = WithUserToken(ctx, token)

	got, err := GetUserTokenFromContext(ctx)
	if err != nil {
		t.Fatalf("expected session error to be cleared, got: %v", err)
	}
	if got != token {
		t.Fatalf("token mismatch\nexpected: %v\nactual:   %v", token, got)
	}
}
