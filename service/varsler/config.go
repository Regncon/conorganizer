package varsler

import (
	"bytes"
	"crypto/ecdh"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"strings"
)

const (
	publicKeyEnvName  = "WEB_PUSH_PUBLIC_KEY"
	privateKeyEnvName = "WEB_PUSH_PRIVATE_KEY"
	subjectEnvName    = "WEB_PUSH_SUBJECT"
)

type Config struct {
	Enabled    bool
	PublicKey  string
	PrivateKey string
	Subject    string
}

func ConfigFromEnv() (Config, error) {
	config := Config{
		PublicKey:  strings.TrimSpace(os.Getenv(publicKeyEnvName)),
		PrivateKey: strings.TrimSpace(os.Getenv(privateKeyEnvName)),
		Subject:    strings.TrimSpace(os.Getenv(subjectEnvName)),
	}
	if config.PublicKey == "" && config.PrivateKey == "" && config.Subject == "" {
		return config, nil
	}
	if config.PublicKey == "" || config.PrivateKey == "" || config.Subject == "" {
		return Config{}, errors.New("WEB_PUSH_PUBLIC_KEY, WEB_PUSH_PRIVATE_KEY, and WEB_PUSH_SUBJECT must all be configured")
	}
	if err := validateVAPIDKeys(config.PublicKey, config.PrivateKey); err != nil {
		return Config{}, err
	}
	if err := validateSubject(config.Subject); err != nil {
		return Config{}, err
	}
	config.Enabled = true
	return config, nil
}

func validateVAPIDKeys(publicKey, privateKey string) error {
	publicBytes, err := base64.RawURLEncoding.DecodeString(publicKey)
	if err != nil || len(publicBytes) != 65 || publicBytes[0] != 4 {
		return errors.New("WEB_PUSH_PUBLIC_KEY must be an uncompressed P-256 public key encoded as base64url")
	}
	curve := ecdh.P256()
	public, err := curve.NewPublicKey(publicBytes)
	if err != nil {
		return errors.New("WEB_PUSH_PUBLIC_KEY is not a valid P-256 point")
	}
	privateBytes, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil || len(privateBytes) != 32 {
		return errors.New("WEB_PUSH_PRIVATE_KEY must be a 32-byte P-256 private key encoded as base64url")
	}
	private, err := curve.NewPrivateKey(privateBytes)
	if err != nil {
		return errors.New("WEB_PUSH_PRIVATE_KEY is outside the valid P-256 range")
	}
	if !bytes.Equal(private.PublicKey().Bytes(), public.Bytes()) {
		return errors.New("WEB_PUSH_PUBLIC_KEY and WEB_PUSH_PRIVATE_KEY do not form a matching pair")
	}
	return nil
}

func validateSubject(subject string) error {
	if strings.HasPrefix(subject, "mailto:") {
		if _, err := mail.ParseAddress(strings.TrimPrefix(subject, "mailto:")); err != nil {
			return fmt.Errorf("WEB_PUSH_SUBJECT contains an invalid mailto address: %w", err)
		}
		return nil
	}
	parsed, err := url.Parse(subject)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return errors.New("WEB_PUSH_SUBJECT must be a mailto address or HTTPS URL")
	}
	return nil
}

func webPushSubscriber(subject string) string {
	return strings.TrimPrefix(subject, "mailto:")
}
