package server

import (
	"net/http"
	"text/template"
)

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("views/authorize.html")

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
	   tmpl, err := template.ParseFiles("views/menu.html", "views/index.html")

	   	if err != nil {
	   		http.Error(w, err.Error(), http.StatusInternalServerError)
	   		return
	   	}

	   	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
	   		http.Error(w, err.Error(), http.StatusInternalServerError)
	   		return
	   	}
	*/
}
