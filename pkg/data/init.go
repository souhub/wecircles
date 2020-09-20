package data

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func init() {
	// Connect to DataBase
	dBUser := os.Getenv("DB_USER")
	dBPass := os.Getenv("DB_PASS")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s", dBUser, dBPass, dbProtocol, dbEndpoint, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Create users table
	cmd := `CREATE TABLE IF NOT EXISTS users(
		id INT AUTO_INCREMENT,
		uuid STRING,
		name STRING,
		userIdStr STRING,
		password STRING,
		image_path STRING,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}

	// Create sessions table
	cmd = `CREATE TABLE IF NOT EXISTS sessions(
		id INT AUTO_INCREMENT
		uuid STRING
		email STRING
		user_id INT
		user_id_str STRING
		created_at TIMESTAMP NOT NULL DEFAULT CURRNET_TIMESTAMP`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}

}

func Encrypt(plainpass string) (cryptpass string) {
	b := []byte(plainpass)
	cryptpass = fmt.Sprintln(sha256.Sum256(b))
	return
}
