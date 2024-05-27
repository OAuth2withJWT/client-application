package api

type BalanceResponse struct {
	UserId     int     `json:"user_id"`
	TotalValue float64 `json:"total_balance"`
}

type AmountResponse struct {
	Category   string  `json:"category"`
	TotalValue float64 `json:"total_amount"`
}
