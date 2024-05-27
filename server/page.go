package server

type Page struct {
	Fields          map[string]string
	Errors          map[string]string
	QueryParameters map[string]string
}
