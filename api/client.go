package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) GetUserBalance(userId int, accessToken string) (BalanceResponse, error) {
	var balance BalanceResponse
	url := "http://localhost:3000/api/cards/balance/" + strconv.Itoa(userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BalanceResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return BalanceResponse{}, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return BalanceResponse{}, err
	}
	json.Unmarshal(body, &balance)

	return balance, nil
}

func (c *Client) GetUserCurrentExpensesByDateAndTime(userId int, accessToken string, date string, time string) (AmountResponse, error) {
	if time == "" {
		time = "00:00:00"
	}

	var amount AmountResponse
	url := "http://localhost:3000/api/transactions/amount/" + strconv.Itoa(userId) + "?date=" + date + "&time=" + time + "&transaction_type=" + "expense"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AmountResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AmountResponse{}, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return AmountResponse{}, err
	}
	json.Unmarshal(body, &amount)

	return amount, nil
}

func (c *Client) GetUserCurrentIncomeByDateAndTime(userId int, accessToken string, date string, time string) (AmountResponse, error) {
	if time == "" {
		time = "00:00:00"
	}

	var amount AmountResponse
	url := "http://localhost:3000/api/transactions/amount/" + strconv.Itoa(userId) + "?date=" + date + "&time=" + time + "&transaction_type=" + "income"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AmountResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AmountResponse{}, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return AmountResponse{}, err
	}
	json.Unmarshal(body, &amount)

	return amount, nil
}

func (c *Client) GetUserCards(userId int, accessToken string) (CardResponse, error) {
	var cardResponse CardResponse
	url := fmt.Sprintf("http://localhost:3000/api/cards/%d", userId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CardResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CardResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return CardResponse{}, fmt.Errorf("failed to fetch cards: status %d", res.StatusCode)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return CardResponse{}, readErr
	}

	err = json.Unmarshal(body, &cardResponse.Cards)
	if err != nil {
		return CardResponse{}, err
	}

	return cardResponse, nil
}

func (c *Client) GetUserTransactions(userId int, accessToken string, date string, time string) (TransactionResponse, error) {
	url := "http://localhost:3000/api/transactions/" + strconv.Itoa(userId) + "?date=" + date

	var transactionResponse TransactionResponse

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TransactionResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TransactionResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return TransactionResponse{}, fmt.Errorf("failed to fetch transactions: status %d", res.StatusCode)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return TransactionResponse{}, readErr
	}

	err = json.Unmarshal(body, &transactionResponse.Transactions)
	if err != nil {
		return TransactionResponse{}, err
	}

	return transactionResponse, nil
}

func (c *Client) GetUserCategoryExpensesByDateAndTime(userId int, accessToken string, date string, time string) (map[string]CategoryExpenseResponse, error) {
	if time == "" {
		time = "00:00:00"
	}

	url := fmt.Sprintf("http://localhost:3000/api/transactions/category_expenses/%d?date=%s&time=%s", userId, date, time)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch category expenses: status %d", res.StatusCode)
	}

	var categoryExpenses []CategoryExpenseResponse
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if err := json.Unmarshal(body, &categoryExpenses); err != nil {
		return nil, err
	}

	categoryExpensesMap := make(map[string]CategoryExpenseResponse)
	for _, expense := range categoryExpenses {
		categoryExpensesMap[expense.ExpenseCategory] = expense
	}

	return categoryExpensesMap, nil
}
