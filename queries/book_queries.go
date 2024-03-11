package queries

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	"github.com/quocquann/locallibrary/models"
)

type IBookQueries interface {
	GetBooks(string, string) ([]models.Book, error)
	AddBooks([]models.Book) error
	GetLastTimeInsertedBook() (string, error)
}

func NewBookQueries(db *sql.DB) IBookQueries {
	return &BookQueries{
		DB: db,
	}
}

type BookQueries struct {
	*sql.DB
}

func (q *BookQueries) GetLastTimeInsertedBook() (string, error) {
	var lastTimeUpdateString string
	query := "SELECT MAX(Updated_at) FROM Book"
	if err := q.QueryRow(query).Scan(&lastTimeUpdateString); err != nil {
		return "", err
	}

	return lastTimeUpdateString, nil
}

func (q *BookQueries) GetBooks(title, authorName string) ([]models.Book, error) {

	books := []models.Book{}

	query := ""
	var rows *sql.Rows
	var err error

	if title == "" && authorName == "" {
		query = "SELECT Book.Id, Book.Isbn, Title, Image, Author.Name, Genre FROM Book JOIN Author ON Book.author_id = Author.id"
		rows, err = q.Query(query)
	} else if title != "" && authorName != "" {
		query = "SELECT Book.Id, Book.Isbn, Title, Image, Author.Name, Genre FROM Book JOIN Author ON Book.author_id = Author.id WHERE Title LIKE CONCAT('%', ?, '%') AND Author.Name LIKE CONCAT('%', ?, '%')"
		rows, err = q.Query(query, title, authorName)
	} else if title != "" {
		query = "SELECT Book.Id, Book.Isbn, Title, Image, Author.Name, Genre FROM Book JOIN Author ON Book.author_id = Author.id WHERE Title LIKE CONCAT('%', ?, '%')"
		rows, err = q.Query(query, title)
	} else {
		query = "SELECT Book.Id, Book.Isbn, Title, Image, Author.Name, Genre FROM Book JOIN Author ON Book.author_id = Author.id WHERE Author.Name LIKE CONCAT('%', ?, '%')"
		rows, err = q.Query(query, authorName)
	}

	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := models.Book{}
		rows.Scan(&book.Id, &book.Isbn, &book.Title, &book.Image, &book.Author.Name, &book.Genre)
		books = append(books, book)
	}

	return books, nil
}

func (q *BookQueries) AddBooks(books []models.Book) error {
	if len(books) <= 0 {
		return nil
	}
	authorQueries := NewAuthorQueries(q.DB)
	if err := authorQueries.AddAuthor(books); err != nil {
		return err
	}
	authorMap, err := authorQueries.GetAuthorMap()
	if err != nil {
		return err
	}

	query := "INSERT INTO Book(Isbn, Title, Image, Author_id, Genre, Updated_at) VALUES "
	vals := []interface{}{}
	b := bytes.Buffer{}
	b.WriteString(query)
	for _, book := range books {
		b.WriteString("(?, ?, ?, ?, ?, ?),")
		vals = append(vals, book.Isbn, book.Title, book.Image, authorMap[book.Author.Name], book.Genre, time.Now())
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Isbn=VALUES(Isbn), Title=VALUES(Title), Image=VALUES(Image), Author_id=VALUES(Author_id), Genre=VALUES(Genre), Updated_at=NOW()")

	query = b.String()

	_, err = q.Exec(query, vals...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
