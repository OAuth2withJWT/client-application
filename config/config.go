package config

import "os"

type IDPConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
}

func LoadIDPConfig() IDPConfig {
	return IDPConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8000/oauth2/callback",
		AuthURL:      "http://localhost:8080/oauth2/auth",
		TokenURL:     "http://localhost:8080/oauth2/token",
	}
}
