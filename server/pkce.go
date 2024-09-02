package server

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
)

const (
	codeVerifierLength = 64
)

func generateCodeVerifier() (string, error) {
	unreserved := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~")

	codeVerifier := make([]byte, codeVerifierLength)
	_, err := rand.Read(codeVerifier)
	if err != nil {
		return "", err
	}

	for i := range codeVerifier {
		codeVerifier[i] = unreserved[int(codeVerifier[i])%len(unreserved)]
	}

	return string(codeVerifier), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
