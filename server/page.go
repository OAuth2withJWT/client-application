package server

import "github.com/OAuth2withJWT/client-application/api"

type Page struct {
	Fields          *map[string]string
	Errors          *map[string]string
	QueryParameters *map[string]string
	Cards           *[]api.Card
	Transactions    *[]api.Transactions
}
