package db

import (
	"database/sql"
	"strconv"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string

	db *sql.DB
)

func getEnv(key string) (string, error) {
    if value, exists := os.LookupEnv(key); exists && value != "" {
        return value, nil
    }
	return "", fmt.Errorf("%s environment variable is missing", key);
}

func getEnvDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func loadConfigs() error {
	var err error

	dbHost, err = getEnv("DB_HOST");
	if (err != nil) {return err}

	dbPort, err = strconv.Atoi(getEnvDefault("DB_PORT", "5432"));
	if (err != nil) {return err}

    dbUser, err = getEnv("DB_USER");
	if (err != nil) {dbUser = "none"}
   
	if (dbUser != "none") {
		dbPassword, err = getEnv("DB_PASSWORD");
		if (err != nil) {dbPassword = "none"}
	}

	dbName, err = getEnv("DB_NAME");
	if (err != nil) {return err}

	return nil
}

func ConnectDatabase() error {
	err := loadConfigs()
	if (err != nil) {
		return err
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	newdb, err := sql.Open("postgres", psqlconn)
	db = newdb
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CheckConnection() error {
	return db.Ping()
}
