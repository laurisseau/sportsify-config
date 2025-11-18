package config

import (
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
)

/*
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
*/

func DB() *sql.DB {

	Secrets := LoadSecretsEnv()

	// Connection parameters
	username := "root"
	password := Secrets["MYSQL_PASSWORD"]// Get from k8s secret
	host := "my-mysql-cluster.sql-server.svc.cluster.local" // K8s service DNS
	port := "6446" // MySQL Router port
	database := ""

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to MySQL!")
	return db

	}