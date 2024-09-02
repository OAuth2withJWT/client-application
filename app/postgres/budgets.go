package postgres

import (
	"database/sql"
	"log"

	"github.com/OAuth2withJWT/client-application/app"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{
		db: db,
	}
}

func (br *BudgetRepository) GetBudgetByUserIdMonthAndCategory(userId int, date string, category string) (app.Budget, error) {
	var budget app.Budget
	err := br.db.QueryRow("SELECT category, amount FROM budgets WHERE user_id = $1 AND category=$2 AND EXTRACT(MONTH FROM month) = EXTRACT(MONTH FROM DATE '"+date+"') AND EXTRACT(YEAR FROM month) = EXTRACT(YEAR FROM DATE '"+date+"');", userId, category).Scan(&budget.Category, &budget.Amount)
	if err != nil {
		return app.Budget{}, err
	}
	return budget, nil
}

func (br *BudgetRepository) GetBudgetsByUserIdAndMonth(userId int, date string) ([]app.Budget, error) {
	rows, err := br.db.Query("SELECT category, amount FROM budgets WHERE user_id = $1 AND EXTRACT(MONTH FROM month) = EXTRACT(MONTH FROM DATE '"+date+"') AND EXTRACT(YEAR FROM month) = EXTRACT(YEAR FROM DATE '"+date+"');", userId)
	if err != nil {
		return []app.Budget{}, err
	}

	var budgets []app.Budget
	for rows.Next() {
		var budget app.Budget
		err := rows.Scan(&budget.Category, &budget.Amount)
		if err != nil {
			log.Fatal(err)
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (br *BudgetRepository) UpdateBudget(userId int, category string, amount float64) error {
	_, err := br.db.Exec(`
		INSERT INTO budgets (user_id, category, amount, month)
		VALUES ($1, $2, $3, CURRENT_DATE)
		ON CONFLICT (user_id, category) 
		DO UPDATE SET amount = EXCLUDED.amount`,
		userId, category, amount,
	)
	return err
}
