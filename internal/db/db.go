package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const (
	tursoDBEnvKey   = "TURSO_DATABASE_URL"
	tursoAuthEnvKey = "TURSO_AUTH_TOKEN"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return db, err
	}

	tursoDBUrl, tursoDBUrlExists := os.LookupEnv(tursoDBEnvKey)
	tursoAuthKey, tursoAuthKeyExists := os.LookupEnv(tursoAuthEnvKey)

	if !tursoDBUrlExists || !tursoAuthKeyExists {
		return db, errors.New("missing turso keys in .env file")
	}

	url := fmt.Sprintf("%v?authToken=%v", tursoDBUrl, tursoAuthKey)

	tursoDB, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}

	db = tursoDB

	return db, err
}

// TODO: Determine if we need this?  The returned db in init may be enough
func GetDB() *sql.DB {
	return db
}
