package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Create NewDatabase
func NewDB() *sql.DB {
	dBUser := os.Getenv("DB_USER")
	dBPass := os.Getenv("DB_PASS")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s", dBUser, dBPass, dbProtocol, dbEndpoint, dbName)
	dbConnection, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return dbConnection

}
