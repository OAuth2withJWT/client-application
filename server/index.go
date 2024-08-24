package server

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/OAuth2withJWT/client-application/app"
)

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {

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

		_, err := s.client.GetUserBalance(userId, accessToken)
		if err != nil {
			log.Print(err.Error())
		}

		now := time.Now()

		beginningOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

		_, err = s.client.GetUserCurrentExpensesByDateAndTime(userId, accessToken, beginningOfMonth, "")
		if err != nil {
			log.Print(err.Error())
		}

		_, err = s.app.BudgetService.GetBudgetsByUserIdAndMonth(userId, now.Format("2006-01-02"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("views/menu.html", "views/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
