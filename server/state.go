package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const (
	stateLength = 32
)

func generateRandomState() (string, error) {
	b := make([]byte, stateLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)
	return state, nil
}

func verifyState(w http.ResponseWriter, r *http.Request, state string) bool {
	cookieState := getStateFromCookie(w, r)
	if cookieState == "" || state == "" {
		return false
	}
	return cookieState == state
}
