package queries

import (
	"bytes"
	"database/sql"

	"github.com/quocquann/locallibrary/models"
)

type BookQueries struct {
	*sql.DB
}

func (q *BookQueries) GetBooks() ([]models.Book, error) {

	books := []models.Book{}

	query := "SELECT Title, Image, Author, Genre FROM Book JOIN Author ON Book.author_id = Author.author_id"
	rows, err := q.Query(query)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := models.Book{}
		rows.Scan(&book.Title, &book.Image, &book.Author, &book.Genre)
		books = append(books, book)
	}

	return books, nil
}

func (q *BookQueries) AddBooks(books []models.Book) error {

	queryString := "INSERT INTO Author(Name) VALUES "
	b := bytes.Buffer{}
	b.WriteString(queryString)
	vals := []interface{}{}
	for i := 0; i < len(books); i++ {
		b.WriteString("(?),")
		vals = append(vals, books[i].Author)
	}

	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Name=VALUES(Name)")

	_, err := q.Exec(b.String(), vals...)
	if err != nil {
		return err
	}

	authorMap := map[string]int{}

	queryString = "SELECT * FROM Author"
	rows, err := q.Query(queryString)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		authorId := 0
		authorName := ""
		rows.Scan(&authorId, &authorName)
		authorMap[authorName] = authorId
	}

	query := "INSERT INTO Book(Title, Image, Author_id, Genre) VALUES "
	vals = []interface{}{}
	b = bytes.Buffer{}
	b.WriteString(query)
	for i := 0; i < len(books); i++ {
		b.WriteString("(?, ?, ?, ?),")
		vals = append(vals, books[i].Title, books[i].Image, authorMap[books[i].Author], books[i].Genre)
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Title=VALUES(Title), Image=VALUES(Image), Author_id=VALUES(Author_id), Genre=VALUES(Genre)")

	query = b.String()

	_, err = q.Exec(query, vals...)
	if err != nil {
		return err
	}

	return nil
}
