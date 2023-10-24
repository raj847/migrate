package config

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var db *sql.DB
var dbCloud *sql.DB

// Open Connection
func OpenConnection() error {
	var err error
	db, err = setupConnection()

	return err
}

func OpenConnectionCloud() error {
	var err error
	dbCloud, err = setupConnectionCloud()

	return err
}

//setupConnection adalah
func setupConnection() (*sql.DB, error) {
	var connection = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		DBUser, DBPass, DBName, DBHost, DBPort, SSLMode)
	fmt.Println("Connection Info:", DBDriver, connection)

	db, err := sql.Open(DBDriver, connection)
	if err != nil {
		return db, errors.New("Connection closed: Failed Connect Database")
	}

	return db, nil
}

func setupConnectionCloud() (*sql.DB, error) {
	var connection = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		DBUserCloud, DBPassCloud, DBNameCloud, DBHostCloud, DBPortCloud, SSLMode)
	fmt.Println("Connection Info:", DBDriver, connection)

	dbCloud, err := sql.Open(DBDriver, connection)
	if err != nil {
		return db, errors.New("Connection closed: Failed Connect Database")
	}

	return dbCloud, nil
}

//CloseConnectionDB adalah
func CloseConnectionDB() {
	db.Close()
}

//DBConnection adalah
func DBConnection() *sql.DB {
	return db
}

//CloseConnectionDB adalah
func CloseConnectionDBCloud() {
	dbCloud.Close()
}

//DBConnection adalah
func DBConnectionCloud() *sql.DB {
	return dbCloud
}
