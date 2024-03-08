package db

import (
	"database/sql"
)

func Init(db *sql.DB) error {
	queryString := `
	CREATE TABLE IF NOT EXISTS Author (
		author_id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(500) UNIQUE
	);`
	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	queryString = `
	CREATE TABLE IF NOT EXISTS Book (
		book_id int PRIMARY KEY AUTO_INCREMENT,
		title varchar(500),
		image varchar(500),
		genre varchar(500),
		author_id int,
		FOREIGN KEY (author_id) REFERENCES Author(author_id)
	);`
	_, err = db.Exec(queryString)
	if err != nil {
		return err
	}

	return nil
}
