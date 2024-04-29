package main

import (
	"log"

	"github.com/OAuth2withJWT/client-application/server"
)

func main() {
	s := server.New()

	log.Fatal(s.Run())
}
