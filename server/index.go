package server

import (
	"net/http"
	"text/template"
)

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	/*
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
	*/

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
