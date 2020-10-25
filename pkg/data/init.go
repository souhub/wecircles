package data

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/souhub/wecircles/pkg/logging"
)

func init() {
	db := NewDB()
	defer db.Close()
	// Create users table
	cmd := `CREATE TABLE IF NOT EXISTS users(
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		user_id_str VARCHAR(255),
		email VARCHAR(255),
		password VARCHAR(255),
		image_path VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`
	_, err := db.Exec(cmd)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Create sessions table
	cmd = `CREATE TABLE IF NOT EXISTS sessions(
		id INT AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(255),
		email VARCHAR(255),
		user_id INT,
		user_id_str VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Create posts table
	cmd = `CREATE TABLE IF NOT EXISTS posts(
			id INT AUTO_INCREMENT PRIMARY KEY,
			uuid VARCHAR(255),
			title VARCHAR(255),
			body TEXT,
			user_id INT,
			user_id_str VARCHAR(255),
			user_name VARCHAR(255),
			thumbnail_path VARCHAR(255),
			created_at  VARCHAR(255)
			)`
	_, err = db.Exec(cmd)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Create  circles table
	cmd = `CREATE TABLE IF NOT EXISTS circles(
			id INT AUTO_INCREMENT PRIMARY KEY,
			Name VARCHAR(50),
			image_path VARCHAR(255),
			overview TEXT,
			category VARCHAR(30),
			owner_id INT,
			owner_id_str VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`
	_, err = db.Exec(cmd)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Create the memberships table
	cmd = `CREATE TABLE IF NOT EXISTS memberships(
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		circle_id INT
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func NewDB() *sql.DB {
	// Connect to DataBase
	var db *sql.DB
	dBUser := os.Getenv("DB_USER")
	dBPass := os.Getenv("DB_PASS")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s", dBUser, dBPass, dbProtocol, dbEndpoint, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	return db
}

func Encrypt(plainpass string) (cryptpass string) {
	b := []byte(plainpass)
	cryptpass = fmt.Sprintln(sha256.Sum256(b))
	return
}
