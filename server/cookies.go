package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"
)

var secretKey []byte

func init() {
	secretKey = []byte(os.Getenv("SECRET_KEY"))
}

func setStateCookie(w http.ResponseWriter, state string) {
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	err := writeSigned(w, cookie, secretKey)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func getStateFromCookie(w http.ResponseWriter, r *http.Request) string {
	value, err := readSigned(r, "oauth_state", secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return ""
	}
	return value
}

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func writeSigned(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	cookie.Value = string(signature) + cookie.Value
	return write(w, cookie)
}

func readSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	signedValue, err := read(r, name)
	if err != nil {
		return "", err
	}

	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return value, nil
}

func write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)

	return nil
}

func read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

func setCodeVerifierCookie(w http.ResponseWriter, codeVerifier string) {
	encodedVerifier := base64.RawURLEncoding.EncodeToString([]byte(codeVerifier))
	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    encodedVerifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getCodeVerifierFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("code_verifier")
	if err != nil {
		return "", err
	}
	decodedVerifier, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", err
	}

	return string(decodedVerifier), nil
}
func deleteCodeVerifierCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
