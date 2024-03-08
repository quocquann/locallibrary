package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func OpenConnection() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USERNAME"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Addr:   os.Getenv("DB_ADDRESS"),
		DBName: os.Getenv("DB_NAME"),
		Net:    "tcp",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}
	fmt.Println("Connected")
	if err := Init(db); err != nil {
		return nil, err
	}
	return db, nil
}
