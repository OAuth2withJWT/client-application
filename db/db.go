package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", getEnvVariable("DATABASE_USER"), getEnvVariable("DATABASE_NAME"), getEnvVariable("DATABASE_PASSWORD"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}

func getEnvVariable(key string) string {
	envVariable, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Could not establish database connection: no %s env var", key)
	}

	return envVariable
}
