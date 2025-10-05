package keyvalue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Regncon/conorganizer/pages/root"
	"github.com/nats-io/nats.go/jetstream"
)

func BroadcastUpdate(kv jetstream.KeyValue, r *http.Request) error {

	ctx := r.Context()
	allKeys, err := kv.Keys(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve keys: %w", err)
	}
	fmt.Printf("All keys in the KeyValue store: %v\n", allKeys)

	for _, sessionID := range allKeys {
		mvc := &root.TodoMVC{}
		fmt.Printf("Processing session ID: %s\n", sessionID)
		fmt.Printf("mvc is: %+v\n", mvc)

		if entry, err := kv.Get(ctx, sessionID); err == nil {
			if err := json.Unmarshal(entry.Value(), mvc); err != nil {
				continue // Ignore unmarshaling errors for other sessions
			}
			mvc.EditingIdx = -1
			if err := saveMVC(ctx, mvc, sessionID, kv); err != nil {
				fmt.Printf("Failed to save MVC for key %s: %v\n", sessionID, err)
			}
		}
	}
	return nil
}

func saveMVC(ctx context.Context, mvc *root.TodoMVC, sessionID string, kv jetstream.KeyValue) error {
	b, err := json.Marshal(mvc)
	if err != nil {
		return fmt.Errorf("failed to marshal mvc: %w", err)
	}
	if _, err := kv.Put(ctx, sessionID, b); err != nil {
		return fmt.Errorf("failed to put key value: %w", err)
	}
	return nil
}
