package db

import (
	"database/sql"
)

func Init(db *sql.DB) error {
	queryString := `
	CREATE TABLE IF NOT EXISTS Author (
		id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(500) UNIQUE
	);`
	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	queryString = `
	CREATE TABLE IF NOT EXISTS Book (
		id int PRIMARY KEY AUTO_INCREMENT,
		isbn varchar(15) UNIQUE,
		title varchar(500),
		image varchar(500),
		genre varchar(500),
		updated_at datetime,
		author_id int,
		FOREIGN KEY (author_id) REFERENCES Author(id),
		INDEX (author_id)
	);`
	_, err = db.Exec(queryString)
	if err != nil {
		return err
	}

	return nil
}
