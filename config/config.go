package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"
)

type IDPConfig struct {
	ClientID               string
	ClientSecret           string
	RedirectURI            string
	AuthURL                string
	TokenURL               string
	IdentityProvider       string
	Scope                  string
	PublicKey              *rsa.PublicKey
	SessionDurationInHours time.Duration
}

func LoadIDPConfig() IDPConfig {
	config := IDPConfig{
		ClientID:               os.Getenv("CLIENT_ID"),
		ClientSecret:           os.Getenv("CLIENT_SECRET"),
		RedirectURI:            "http://localhost:8000/oauth2/callback",
		AuthURL:                "http://localhost:8080/oauth2/auth",
		TokenURL:               "http://localhost:8080/oauth2/token",
		IdentityProvider:       "http://localhost:8080",
		Scope:                  "transactions:read,cards:read,openid",
		SessionDurationInHours: 24 * time.Hour,
	}

	publicKey, err := LoadPublicKey("keys/public_key.pem")
	if err != nil {
		log.Printf("Failed to load public key, continuing with default config: %v", err)
	} else {
		config.PublicKey = publicKey
	}

	return config
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pubKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}

	return rsaPublicKey, nil
}
