package server

import (
	"fmt"
	"net/http"
)

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, "Error generating state", http.StatusInternalServerError)
		return
	}
	setStateCookie(w, state)

	http.Redirect(w, r, fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s", s.IDPConfig.AuthURL, s.IDPConfig.ClientID, s.IDPConfig.RedirectURI, state), http.StatusFound)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state == "" || code == "" {
		http.Error(w, "Missing state or code in the callback", http.StatusBadRequest)
		return
	}

	if verifyState(r, state) {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Authorization code: %s", code)

}
