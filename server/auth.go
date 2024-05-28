package server

import (
	"fmt"
	"net/http"
)

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	state := generateState()

	setStateCookie(w, state)

	http.Redirect(w, r, fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s", s.IDPConfig.AuthURL, s.IDPConfig.ClientID, s.IDPConfig.RedirectURI, state), http.StatusFound)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	receivedState := r.URL.Query().Get("state")

	storedState := getStateFromCookie(r)

	if receivedState == storedState {
		fmt.Fprintf(w, "Authorization code: %s", code)
	} else {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
	}
}
