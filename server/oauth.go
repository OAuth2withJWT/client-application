package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OAuth2withJWT/client-application/app"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

var userId = 1

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, "Error generating state", http.StatusInternalServerError)
		return
	}
	setStateCookie(w, state)

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		http.Error(w, "Error generating code verifier", http.StatusInternalServerError)
		return
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	setCodeVerifierCookie(w, codeVerifier)

	http.Redirect(w, r, fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=S256", s.IDPConfig.AuthURL, s.IDPConfig.ClientID, s.IDPConfig.RedirectURI, state, codeChallenge), http.StatusFound)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state == "" || code == "" {
		http.Error(w, "Missing state or code in the callback", http.StatusBadRequest)
		return
	}

	if !verifyState(w, r, state) {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
		return
	}

	codeVerifier, err := getCodeVerifierFromCookie(r)
	if err != nil {
		http.Error(w, "Failed to get code verifier", http.StatusInternalServerError)
		return
	}

	deleteCodeVerifierCookie(w)

	tokenReq := TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		ClientID:     s.IDPConfig.ClientID,
		ClientSecret: s.IDPConfig.ClientSecret,
		RedirectURI:  s.IDPConfig.RedirectURI,
		CodeVerifier: codeVerifier,
	}

	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(tokenReq); err != nil {
		http.Error(w, "Failed to encode token request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(s.IDPConfig.TokenURL, "application/json", &reqBody)
	if err != nil {
		http.Error(w, "Failed to make token request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Token request failed: %s", string(bodyBytes)), http.StatusInternalServerError)
		return
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		return
	}

	sessionID, err := s.app.SessionService.CreateSession(userId, tokenResp.AccessToken, time.Now().Add(app.SessionDurationInHours*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setAccessCookie(w, sessionID)
	http.Redirect(w, r, "/", http.StatusFound)
}
