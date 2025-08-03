package config

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
)

// DB establishes a connection to the PostgreSQL database using secrets from AWS
func DB() *sql.DB {

	Secrets := LoadSecretsEnv()

	// Access secrets
	host := Secrets["POSTGRES_ENDPOINT"]
	port := Secrets["POSTGRES_PORT"]
	user := Secrets["POSTGRES_USERNAME"]
	password := Secrets["POSTGRES_PASSWORD"]
	dbname := Secrets["POSTGRES_DATABASE"]

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	)

	// Connect to DB
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	fmt.Println("✅ Connected to DB")
	return db
}
