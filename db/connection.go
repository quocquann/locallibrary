package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/quocquann/locallibrary/queries"
)

type Queries struct {
	*queries.BookQueries
}

func OpenConnection() (*Queries, error) {
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
	return &Queries{
		BookQueries: &queries.BookQueries{DB: db},
	}, nil
}
