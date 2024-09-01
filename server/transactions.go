package server

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/OAuth2withJWT/client-application/app"
)

func (s *Server) handleTransactionsPage(w http.ResponseWriter, r *http.Request) {

	sessionID := getAuthSessionIDFromCookie(r)
	session, err := s.app.SessionService.ValidateSession(sessionID)

	if (err != nil || session == app.Session{}) {
		deleteAuthSessionCookie(w)
		tmpl, _ := template.ParseFiles("views/authorize.html")

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		accessToken := session.AccessToken
		user, _ := s.GetUserInfoFromIDToken(session.IdToken)
		userId := user.ID

		today := time.Now().Format("2006-01-02")

		transactionResponse, err := s.client.GetUserTransactions(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		page := Page{
			Transactions: &transactionResponse.Transactions,
		}

		tmpl, err := template.ParseFiles("views/menu.html", "views/transactions.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "transactions.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
