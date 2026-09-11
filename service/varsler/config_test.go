package varsler

import (
	"crypto/ecdh"
	"encoding/base64"
	"testing"
)

func TestConfigFromEnv_WhenUnsetDisablesWebPush(t *testing.T) {
	// Given
	t.Setenv("WEB_PUSH_PUBLIC_KEY", "")
	t.Setenv("WEB_PUSH_PRIVATE_KEY", "")
	t.Setenv("WEB_PUSH_SUBJECT", "")

	// When
	config, err := ConfigFromEnv()

	// Then
	if err != nil {
		t.Fatalf("expected empty configuration to be accepted: %v", err)
	}
	if config.Enabled {
		t.Fatal("expected web push to be disabled")
	}
}

func TestConfigFromEnv_WhenPartiallyConfiguredReturnsError(t *testing.T) {
	// Given
	t.Setenv("WEB_PUSH_PUBLIC_KEY", "public")
	t.Setenv("WEB_PUSH_PRIVATE_KEY", "")
	t.Setenv("WEB_PUSH_SUBJECT", "mailto:varsler@example.com")

	// When
	_, err := ConfigFromEnv()

	// Then
	if err == nil {
		t.Fatal("expected partial web push configuration to fail")
	}
}

func TestConfigFromEnv_WhenKeysAreMalformedReturnsError(t *testing.T) {
	// Given
	t.Setenv("WEB_PUSH_PUBLIC_KEY", "not-a-key")
	t.Setenv("WEB_PUSH_PRIVATE_KEY", "not-a-key")
	t.Setenv("WEB_PUSH_SUBJECT", "mailto:varsler@example.com")

	// When
	_, err := ConfigFromEnv()

	// Then
	if err == nil {
		t.Fatal("expected malformed VAPID keys to fail")
	}
}

func TestConfigFromEnv_WhenKeyPairDoesNotMatchReturnsError(t *testing.T) {
	// Given
	privateOne := make([]byte, 32)
	privateOne[31] = 1
	privateTwo := make([]byte, 32)
	privateTwo[31] = 2
	private, err := ecdh.P256().NewPrivateKey(privateOne)
	if err != nil {
		t.Fatalf("create test key: %v", err)
	}
	t.Setenv("WEB_PUSH_PUBLIC_KEY", base64.RawURLEncoding.EncodeToString(private.PublicKey().Bytes()))
	t.Setenv("WEB_PUSH_PRIVATE_KEY", base64.RawURLEncoding.EncodeToString(privateTwo))
	t.Setenv("WEB_PUSH_SUBJECT", "mailto:varsler@example.com")

	// When
	_, err = ConfigFromEnv()

	// Then
	if err == nil {
		t.Fatal("expected mismatched VAPID key pair to fail")
	}
}

func TestWebPushSubscriber_StripsMailtoPrefixExpectedByDependency(t *testing.T) {
	// Given
	subject := "mailto:varsler@example.com"

	// When
	subscriber := webPushSubscriber(subject)

	// Then
	if subscriber != "varsler@example.com" {
		t.Fatalf("expected bare email for Web Push dependency, got %q", subscriber)
	}
}
