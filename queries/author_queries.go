package queries

import (
	"bytes"
	"database/sql"

	"github.com/quocquann/locallibrary/models"
)

type IAuthorQueries interface {
	AddAuthor([]models.Book) error
	GetAuthorMap() (map[string]int, error)
}

func NewAuthorQueries(db *sql.DB) IAuthorQueries {
	return &AuthorQueries{
		DB: db,
	}
}

type AuthorQueries struct {
	*sql.DB
}

func (q *AuthorQueries) AddAuthor(books []models.Book) error {
	if len(books) <= 0 {
		return nil
	}
	queryString := "INSERT INTO library_author(Name) VALUES "
	b := bytes.Buffer{}
	b.WriteString(queryString)
	vals := []interface{}{}
	for _, book := range books {
		b.WriteString("(?),")
		vals = append(vals, book.Author.Name)
	}

	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Name=VALUES(Name)")

	_, err := q.Exec(b.String(), vals...)
	if err != nil {
		return err
	}
	return nil
}

func (q *AuthorQueries) GetAuthorMap() (map[string]int, error) {
	authorMap := map[string]int{}

	queryString := "SELECT Id, Name FROM library_author"
	rows, err := q.Query(queryString)
	if err != nil {
		return authorMap, err
	}
	defer rows.Close()
	for rows.Next() {
		authorId := 0
		authorName := ""
		rows.Scan(&authorId, &authorName)
		authorMap[authorName] = authorId
	}
	return authorMap, nil
}
