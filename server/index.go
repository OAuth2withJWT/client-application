package server

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	userId := 1

	balance, err := s.client.GetUserBalance(userId)
	if err != nil {
		log.Print(err.Error())
	}

	now := time.Now()
	beginningOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	amount, err := s.client.GetUserCurrentExpensesByDateAndTime(userId, beginningOfMonth, "")
	if err != nil {
		log.Print(err.Error())
	}

	budgets, err := s.app.BudgetService.GetBudgetsByUserIdAndMonth(userId, now.Format("2006-01-02"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := Page{
		Fields: map[string]string{
			"CurrentMonth":     now.Month().String(),
			"Username":         "Lejla",
			"Balance":          fmt.Sprintf("%.2f", balance.TotalValue),
			"MonthlyBudget":    fmt.Sprintf("%.2f", budgets["monthly"].Amount),
			"Expenses":         fmt.Sprintf("%.2f", amount.TotalValue),
			"HealthcareBudget": fmt.Sprintf("%.2f", budgets["healthcare"].Amount),
		},
	}

	tmpl, _ := template.ParseFiles("views/index.html")
	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
