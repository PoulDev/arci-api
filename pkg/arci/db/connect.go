package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	db *sql.DB
)

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value, nil
	}
	return "", fmt.Errorf("%s environment variable is missing", key)
}

func getEnvDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func loadConfigs() error {
	var err error

	dbHost, err = getEnv("DB_HOST")
	if err != nil {
		return err
	}

	dbPort = getEnvDefault("DB_PORT", "3306")

	dbUser, err = getEnv("DB_USER")
	if err != nil {
		return err
	}

	dbPassword = getEnvDefault("DB_PASSWORD", "")

	dbName, err = getEnv("DB_NAME")
	if err != nil {
		return err
	}

	return nil
}

func ConnectDatabase() error {
	err := loadConfigs()
	if err != nil {
		return err
	}

	// DSN format: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	newdb, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db = newdb

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CheckConnection() error {
	return db.Ping()
}
