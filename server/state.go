package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(b)
	return state, nil
}

func verifyState(r *http.Request, state string) bool {
	cookieState := getStateFromCookie(r)
	if cookieState == "" || state == "" {
		return false
	}
	return cookieState == state
}
