package main

import (
	"log"

	"github.com/OAuth2withJWT/client-application/api"
	"github.com/OAuth2withJWT/client-application/app"
	"github.com/OAuth2withJWT/client-application/app/postgres"
	"github.com/OAuth2withJWT/client-application/db"
	"github.com/OAuth2withJWT/client-application/server"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	budgetRepository := postgres.NewBudgetRepository(db)

	app := app.Application{
		BudgetService: app.NewBudgetService(budgetRepository),
	}

	api := api.Client{}

	s := server.New(&app, &api)

	log.Fatal(s.Run())
}
