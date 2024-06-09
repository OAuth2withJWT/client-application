package server

import (
	"net/http"
	"text/template"
)

func (s *Server) handleBudgetsPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("views/menu.html", "views/budgets.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "budgets.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
