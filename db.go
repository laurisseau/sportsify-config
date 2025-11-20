package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func DB() (*sql.DB, error) {

	Secrets := LoadSecretsEnv()

	// Connection parameters
	username := Secrets["MYSQL_USERNAME"]
	password := Secrets["MYSQL_PASSWORD"]
	host := Secrets["MYSQL_HOST"]
	port := Secrets["MYSQL_PORT"]
	database := Secrets["MYSQL_DATABASE"]

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to MySQL!")
	return db, nil

}
