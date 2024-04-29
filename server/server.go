package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	s := &Server{
		router: mux.NewRouter(),
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
	s.router.HandleFunc("/", s.handleHomePage).Methods("GET")
}
