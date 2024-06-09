package server

import (
	"net/http"
	"text/template"
)

func (s *Server) handleTransactionsPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("views/menu.html", "views/transactions.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "transactions.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
