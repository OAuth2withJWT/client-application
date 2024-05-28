package server

import "net/http"

func setStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getStateFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return ""
	}
	return cookie.Value
}
