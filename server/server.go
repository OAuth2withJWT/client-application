package server

import (
	"log"
	"net/http"

	"github.com/OAuth2withJWT/client-application/api"
	"github.com/OAuth2withJWT/client-application/app"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	app    *app.Application
	client *api.Client
}

func New(a *app.Application, c *api.Client) *Server {
	s := &Server{
		router: mux.NewRouter(),
		app:    a,
		client: c,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	log.Println("Server started on port 8000")
	return http.ListenAndServe(":8020", s.router)
}

func (s *Server) setupRoutes() {
	s.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	s.router.HandleFunc("/", s.handleIndexPage).Methods("GET")
	s.router.HandleFunc("/transactions", s.handleTransactionsPage).Methods("GET")
	s.router.HandleFunc("/budgets", s.handleBudgetsPage).Methods("GET")
}
