package api

type BalanceResponse struct {
	UserId     int     `json:"user_id"`
	TotalValue float64 `json:"total_balance"`
}

type AmountResponse struct {
	Category   string  `json:"category"`
	TotalValue float64 `json:"total_amount"`
}

type Card struct {
	Id             int     `json:"card_id"`
	UserId         int     `json:"user_id"`
	CardNumber     string  `json:"card_number"`
	CurrentBalance float64 `json:"current_balance"`
	ExpirationDate string  `json:"expiration_date"`
	CardType       string  `json:"card_type"`
}

type CardResponse struct {
	Cards []Card `json:"cards"`
}

type Transactions struct {
	Id                   int     `json:"id"`
	CardId               int     `json:"card_id"`
	Time                 string  `json:"time"`
	Amount               float64 `json:"amount"`
	ExpenseCategory      string  `json:"expense_category"`
	TransactionType      string  `json:"transaction_type"`
	Location             *string `json:"location"`
	DestinationAccountId *int    `json:"destination_account_id"`
	SourceAccountId      *int    `json:"source_account_id"`
}

type TransactionResponse struct {
	Transactions []Transactions `json:"transactions"`
}

type CategoryExpenseResponse struct {
	ExpenseCategory string  `json:"expense_category"`
	TotalSpent      float64 `json:"total_spent"`
}
