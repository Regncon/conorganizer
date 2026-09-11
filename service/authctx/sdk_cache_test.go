package authctx

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestSessionValidator_ReusesCachedSigningKeyWhenKeyServiceIsUnavailable(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A shared Descope validator has fetched a signing key after a retryable failure.",
		When:  "The key service becomes unavailable and another session is validated.",
		Then:  "The cached key validates the session without another network request.",
	})

	// Given
	const projectID = "cache-test-project"
	sessionToken, publicKey := signedSessionTokenAndPublicKey(t, projectID)
	var keyRequests atomic.Int64
	var keyServiceAvailable atomic.Bool
	keyServiceAvailable.Store(true)

	keyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/keys/"+projectID {
			http.NotFound(w, r)
			return
		}

		requestNumber := keyRequests.Add(1)
		if requestNumber == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if !keyServiceAvailable.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"keys":[%s]}`, publicKey)
	}))
	t.Cleanup(keyServer.Close)

	t.Setenv(descope.EnvironmentVariablePublicKey, "")
	t.Setenv(descope.EnvironmentVariableManagementKey, "")
	t.Setenv(descope.EnvironmentVariableAuthManagementKey, "")
	validator, err := newSessionValidator(&client.Config{
		ProjectID:      projectID,
		DescopeBaseURL: keyServer.URL,
		RequestTimeout: time.Second,
	}, discardLogger())
	if err != nil {
		t.Fatalf("create session validator: %v", err)
	}

	valid, _, err := validator.ValidateSessionWithToken(context.Background(), sessionToken)
	if err != nil || !valid {
		t.Fatalf("warm signing-key cache: valid=%v error=%v", valid, err)
	}
	if keyRequests.Load() != 2 {
		t.Fatalf("expected one retry followed by a successful key request, got %d requests", keyRequests.Load())
	}

	// When
	keyServiceAvailable.Store(false)
	valid, _, err = validator.ValidateSessionWithToken(context.Background(), sessionToken)

	// Then
	if err != nil || !valid {
		t.Fatalf("validate session from cached signing key: valid=%v error=%v", valid, err)
	}
	if keyRequests.Load() != 2 {
		t.Fatalf("cached validation made another key request; got %d requests", keyRequests.Load())
	}
}

func signedSessionTokenAndPublicKey(t *testing.T, projectID string) (string, string) {
	t.Helper()

	const keyID = "cache-test-key"
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate signing key: %v", err)
	}
	publicKey, err := jwk.FromRaw(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("create public JWK: %v", err)
	}
	if err := publicKey.Set(jwk.KeyIDKey, keyID); err != nil {
		t.Fatalf("set public key ID: %v", err)
	}
	if err := publicKey.Set(jwk.AlgorithmKey, jwa.ES256); err != nil {
		t.Fatalf("set public key algorithm: %v", err)
	}
	if err := publicKey.Set(jwk.KeyUsageKey, "sig"); err != nil {
		t.Fatalf("set public key usage: %v", err)
	}
	publicKeyJSON, err := json.Marshal(publicKey)
	if err != nil {
		t.Fatalf("marshal public JWK: %v", err)
	}

	token := jwt.New()
	claims := map[string]any{
		jwt.AudienceKey:   []string{projectID},
		jwt.SubjectKey:    "cache-test-user",
		jwt.IssuedAtKey:   time.Now(),
		jwt.ExpirationKey: time.Now().Add(time.Hour),
		jwt.IssuerKey:     projectID,
		"drn":             "DS",
	}
	for name, value := range claims {
		if err := token.Set(name, value); err != nil {
			t.Fatalf("set token claim %q: %v", name, err)
		}
	}

	headers := jws.NewHeaders()
	if err := headers.Set(jws.KeyIDKey, keyID); err != nil {
		t.Fatalf("set token key ID: %v", err)
	}
	if err := headers.Set(jws.TypeKey, "JWT"); err != nil {
		t.Fatalf("set token type: %v", err)
	}
	if err := headers.Set(jws.AlgorithmKey, jwa.ES256); err != nil {
		t.Fatalf("set token algorithm: %v", err)
	}
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.ES256, privateKey, jws.WithProtectedHeaders(headers)))
	if err != nil {
		t.Fatalf("sign session token: %v", err)
	}

	return string(signedToken), string(publicKeyJSON)
}
