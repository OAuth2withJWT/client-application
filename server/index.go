package server

import (
	"fmt"
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

		balance, err := s.client.GetUserBalance(userId, accessToken)
		if err != nil {
			log.Print(err.Error())
		}

		today := time.Now().Format("2006-01-02")

		todaySpending, err := s.client.GetUserCurrentExpensesByDateAndTime(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		todayIncome, err := s.client.GetUserCurrentIncomeByDateAndTime(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		budgets, err := s.app.BudgetService.GetBudgetsByUserIdAndMonth(userId, today)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cardResponse, err := s.client.GetUserCards(userId, accessToken)
		if err != nil {
			log.Print(err.Error())
		}

		transactionResponse, err := s.client.GetUserTransactions(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		totalExpenses := budgets["healthcare"].Amount +
			budgets["groceries"].Amount +
			budgets["clothing"].Amount +
			budgets["entertainment"].Amount +
			budgets["dining"].Amount +
			budgets["transport"].Amount +
			budgets["utilities"].Amount

		monthlyBudget := budgets["monthly"].Amount
		remainingBudget := monthlyBudget - totalExpenses
		remainingBudgetPercentage := (remainingBudget / monthlyBudget) * 100

		page := Page{
			Fields: &map[string]string{
				"Username":                  user.Name,
				"Balance":                   fmt.Sprintf("%.2f", balance.TotalValue),
				"TodaySpending":             fmt.Sprintf("%.2f", todaySpending.TotalValue),
				"TodayIncome":               fmt.Sprintf("%.2f", todayIncome.TotalValue),
				"MonthlyBudget":             fmt.Sprintf("%.0f", budgets["monthly"].Amount),
				"HealthcareBudget":          fmt.Sprintf("%.0f", budgets["healthcare"].Amount),
				"GroceriesBudget":           fmt.Sprintf("%.0f", budgets["groceries"].Amount),
				"TransportBudget":           fmt.Sprintf("%.0f", budgets["transport"].Amount),
				"ClothingBudget":            fmt.Sprintf("%.0f", budgets["clothing"].Amount),
				"EntertainmentBudget":       fmt.Sprintf("%.0f", budgets["entertainment"].Amount),
				"DiningBudget":              fmt.Sprintf("%.0f", budgets["dining"].Amount),
				"UtilitiesBudget":           fmt.Sprintf("%.0f", budgets["utilities"].Amount),
				"RemainingBudgetPercentage": fmt.Sprintf("%.2f", remainingBudgetPercentage),
			},
			Cards:        &cardResponse.Cards,
			Transactions: &transactionResponse.Transactions,
		}

		tmpl, err := template.ParseFiles("views/menu.html", "views/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
