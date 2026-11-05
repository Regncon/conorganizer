package main

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestListenAddressUsesEveryInterfaceWhenHostIsUnset(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given HOST is not set, as in local development and Docker.",
		When:  "When the server address is built.",
		Then:  "Then the server listens on every interface.",
	})

	// Given
	expectedAddress := "0.0.0.0:8080"

	// When
	address := listenAddress("", "8080")

	// Then
	if address != expectedAddress {
		t.Fatalf("expected %q, got %q", expectedAddress, address)
	}
}

func TestListenAddressUsesLoopbackWhenHostIsLoopback(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given HOST is 127.0.0.1, as in the production systemd units.",
		When:  "When the server address is built.",
		Then:  "Then the server only listens on loopback, so only Caddy on the same machine can reach it.",
	})

	// Given
	expectedAddress := "127.0.0.1:19080"

	// When
	address := listenAddress("127.0.0.1", "19080")

	// Then
	if address != expectedAddress {
		t.Fatalf("expected %q, got %q", expectedAddress, address)
	}
}
