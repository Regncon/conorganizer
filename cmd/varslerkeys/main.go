package main

import (
	"fmt"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate VAPID keys: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("WEB_PUSH_PUBLIC_KEY=%s\nWEB_PUSH_PRIVATE_KEY=%s\n", publicKey, privateKey)
}
