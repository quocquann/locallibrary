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

	query := "SELECT * FROM Books"
	rows, err := q.Query(query)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := models.Book{}
		rows.Scan(&book.Isbn, &book.Title, &book.Image, &book.Author, &book.Genre)
		books = append(books, book)
	}

	return books, nil
}

func (q *BookQueries) AddBooks(books []models.Book) error {
	query := "INSERT INTO Books(Isbn, Title, Image, Author, Genre) VALUES "
	vals := []interface{}{}
	b := bytes.Buffer{}
	b.WriteString(query)
	for i := 0; i < len(books); i++ {
		b.WriteString("(?, ?, ?, ?, ?),")
		vals = append(vals, books[i].Isbn, books[i].Title, books[i].Image, books[i].Author, books[i].Genre)
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Isbn=VALUES(Isbn), Title=VALUES(Title), Image=VALUES(Image), Author=VALUES(Author), Genre=VALUES(Genre)")

	query = b.String()

	_, err := q.Exec(query, vals...)
	if err != nil {
		return err
	}

	return nil
}
