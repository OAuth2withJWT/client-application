package app

import (
	"time"
)

type BudgetService struct {
	repository BudgetRepository
}

func NewBudgetService(br BudgetRepository) *BudgetService {
	return &BudgetService{
		repository: br,
	}
}

type BudgetRepository interface {
	GetBudgetByUserIdMonthAndCategory(userId int, date string, category string) (Budget, error)
	GetBudgetsByUserIdAndMonth(userId int, date string) ([]Budget, error)
	UpdateBudget(userId int, category string, amount float64) error
}

type Budget struct {
	Category    string
	Amount      float64
	UpdateStamp time.Time
}

func (s *BudgetService) GetBudgetByUserIdMonthAndCategory(userId int, date string, category string) (Budget, error) {
	budget, err := s.repository.GetBudgetByUserIdMonthAndCategory(userId, date, category)
	if err != nil {
		return Budget{}, err
	}
	return budget, nil
}

func (s *BudgetService) GetBudgetsByUserIdAndMonth(userId int, date string) (map[string]Budget, error) {
	budgets, err := s.repository.GetBudgetsByUserIdAndMonth(userId, date)
	if err != nil {
		return map[string]Budget{}, err
	}

	budgetMap := make(map[string]Budget)
	for _, budget := range budgets {
		budgetMap[budget.Category] = budget
	}
	return budgetMap, nil
}

func (s *BudgetService) UpdateBudget(userId int, category string, amount float64) error {
	return s.repository.UpdateBudget(userId, category, amount)
}
