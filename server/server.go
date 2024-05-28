package server

import (
	"log"
	"net/http"

	"github.com/OAuth2withJWT/client-application/api"
	"github.com/OAuth2withJWT/client-application/app"
	"github.com/OAuth2withJWT/client-application/config"
	"github.com/gorilla/mux"
)

type IDPConfig struct {
	ClientID    string
	RedirectURI string
}

type Server struct {
	router    *mux.Router
	app       *app.Application
	client    *api.Client
	IDPConfig config.IDPConfig
}

func New(a *app.Application, c *api.Client) *Server {
	s := &Server{
		router:    mux.NewRouter(),
		app:       a,
		client:    c,
		IDPConfig: config.LoadIDPConfig(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	log.Println("Server started on port 8000")
	return http.ListenAndServe(":8000", s.router)
}

func (s *Server) setupRoutes() {
	s.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	s.router.HandleFunc("/", s.handleIndexPage).Methods("GET")
	s.router.HandleFunc("/oauth2/callback", s.handleCallback).Methods("GET")
	s.router.HandleFunc("/oauth2/auth", s.handleAuth).Methods("GET")
}
