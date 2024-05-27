package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) GetUserBalance(userId int) (BalanceResponse, error) {
	var balance BalanceResponse
	url := "http://localhost:3000/api/cards/balance/" + strconv.Itoa(userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BalanceResponse{}, err
	}

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

func (c *Client) GetUserCurrentExpensesByDateAndTime(userId int, date string, time string) (AmountResponse, error) {
	if time == "" {
		time = "00:00:00"
	}

	var amount AmountResponse
	url := "http://localhost:3000/api/transactions/amount/" + strconv.Itoa(userId) + "?date=" + date + "&time=" + time

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AmountResponse{}, err
	}

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
