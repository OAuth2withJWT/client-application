package server

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
)

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)
	log.Printf("Generated State: '%s'", state)
	return state, nil
}

func verifyState(r *http.Request, state string) bool {
	cookieState := getStateFromCookie(r)
	log.Print(cookieState)
	log.Print(state)
	if cookieState == "" || state == "" {
		return false
	}
	return cookieState == state
}
