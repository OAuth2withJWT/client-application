package server

import (
	"net/http"
	"text/template"
)

func (s *Server) handleHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/index.html")
	err := tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
