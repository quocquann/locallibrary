package queries

import (
	"bytes"
	"database/sql"
	"fmt"

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
	query := "SELECT MAX(Updated_at) FROM library_book"
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

	//==========ADD AUTHOR============
	authorQueries := NewAuthorQueries(q.DB)
	if err := authorQueries.AddAuthor(books); err != nil {
		return err
	}
	authorMap, err := authorQueries.GetAuthorMap()
	if err != nil {
		return err
	}

	//=========ADD GENRE=============
	genreQueries := NewGenreQueries(q.DB)
	if err := genreQueries.AddGenre(books); err != nil {
		return err
	}
	genreMap, err := genreQueries.GetGenreMap()
	if err != nil {
		return err
	}

	publisherQueries := NewPublisherQueries(q.DB)
	if err := publisherQueries.AddPublisher(books); err != nil {
		return err
	}
	publisherMap, err := publisherQueries.GetPublisherMap()
	if err != nil {
		return err
	}

	query := "INSERT INTO library_book(Isbn, Title, Image, `Describe`, Author_id, Genre_id, Publisher_id) VALUES "
	vals := []interface{}{}
	b := bytes.Buffer{}
	b.WriteString(query)
	for _, book := range books {
		b.WriteString("(?, ?, ?, ?, ?, ?, ?),")
		vals = append(vals, book.Isbn, book.Title, book.Image, book.Describe, authorMap[book.Author.Name], genreMap[book.Genre.Name], publisherMap[book.Publisher.Name])
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Isbn=VALUES(Isbn), Title=VALUES(Title), Image=VALUES(Image), `Describe`=VALUES(`Describe`), Author_id=VALUES(Author_id), Genre_id=VALUES(Genre_id), Publisher_id=VALUES(Publisher_id)")

	query = b.String()

	_, err = q.Exec(query, vals...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
