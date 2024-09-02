package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/OAuth2withJWT/client-application/app"
)

func (s *Server) handleBudgetsPage(w http.ResponseWriter, r *http.Request) {

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

		todayExpenses, err := s.client.GetUserCurrentExpensesByDateAndTime(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		todayIncomes, err := s.client.GetUserCurrentIncomeByDateAndTime(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		budgets, err := s.app.BudgetService.GetBudgetsByUserIdAndMonth(userId, today)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		categoryExpenses, err := s.client.GetUserCategoryExpensesByDateAndTime(userId, accessToken, today, "")
		if err != nil {
			log.Print(err.Error())
		}

		page := Page{
			Fields: &map[string]string{
				"Balance":                      fmt.Sprintf("%.2f", balance.TotalValue),
				"TodayExpenses":                fmt.Sprintf("%.2f", todayExpenses.TotalValue),
				"TodayIncome":                  fmt.Sprintf("%.2f", todayIncomes.TotalValue),
				"MonthlyBudget":                fmt.Sprintf("%.2f", budgets["monthly"].Amount),
				"GroceriesBudget":              fmt.Sprintf("%.2f", budgets["groceries"].Amount),
				"HealthcareBudget":             fmt.Sprintf("%.2f", budgets["healthcare"].Amount),
				"TransportBudget":              fmt.Sprintf("%.2f", budgets["transport"].Amount),
				"ClothingBudget":               fmt.Sprintf("%.2f", budgets["clothing"].Amount),
				"EntertainmentBudget":          fmt.Sprintf("%.2f", budgets["entertainment"].Amount),
				"DiningBudget":                 fmt.Sprintf("%.2f", budgets["dining"].Amount),
				"UtilitiesBudget":              fmt.Sprintf("%.2f", budgets["utilities"].Amount),
				"RemainingBudget":              fmt.Sprintf("%.2f", remainingBudget),
				"GroceriesBudgetRemaining":     fmt.Sprintf("%.2f", budgets["groceries"].Amount-categoryExpenses["groceries"].TotalSpent),
				"HealthcareBudgetRemaining":    fmt.Sprintf("%.2f", budgets["healthcare"].Amount-categoryExpenses["healthcare"].TotalSpent),
				"TransportBudgetRemaining":     fmt.Sprintf("%.2f", budgets["transport"].Amount-categoryExpenses["transport"].TotalSpent),
				"ClothingBudgetRemaining":      fmt.Sprintf("%.2f", budgets["clothing"].Amount-categoryExpenses["clothing"].TotalSpent),
				"EntertainmentBudgetRemaining": fmt.Sprintf("%.2f", budgets["entertainment"].Amount-categoryExpenses["entertainment"].TotalSpent),
				"DiningBudgetRemaining":        fmt.Sprintf("%.2f", budgets["dining"].Amount-categoryExpenses["dining"].TotalSpent),
				"UtilitiesBudgetRemaining":     fmt.Sprintf("%.2f", budgets["utilities"].Amount-categoryExpenses["utilities"].TotalSpent),
			},
		}

		tmpl, err := template.ParseFiles("views/menu.html", "views/budgets.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "budgets.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func (s *Server) handleUpdateBudget(w http.ResponseWriter, r *http.Request) {
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
		user, _ := s.GetUserInfoFromIDToken(session.IdToken)
		userId := user.ID

		var request struct {
			Category     string `json:"category"`
			BudgetAmount string `json:"budgetAmount"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		budgetAmount, err := strconv.ParseFloat(request.BudgetAmount, 64)
		if err != nil {
			http.Error(w, "Invalid budget amount", http.StatusBadRequest)
			return
		}

		err = s.app.BudgetService.UpdateBudget(userId, request.Category, budgetAmount)
		if err != nil {
			http.Error(w, "Failed to update budget", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Budget updated successfully"})
	}
}
