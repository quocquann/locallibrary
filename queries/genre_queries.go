package queries

import (
	"bytes"
	"database/sql"

	"github.com/quocquann/locallibrary/models"
)

type IGenreQuries interface {
	AddGenre([]models.Book) error
	GetGenreMap() (map[string]int, error)
}

func NewGenreQueries(db *sql.DB) IGenreQuries {
	return &GenreQueries{
		DB: db,
	}
}

type GenreQueries struct {
	*sql.DB
}

func (q *GenreQueries) AddGenre(books []models.Book) error {
	if len(books) <= 0 {
		return nil
	}

	query := "INSERT INTO library_genre(Name) VALUES "
	b := bytes.Buffer{}
	vals := []interface{}{}
	b.WriteString(query)
	for _, book := range books {
		b.WriteString("(?),")
		vals = append(vals, book.Genre.Name)
	}

	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Name=VALUES(Name)")

	_, err := q.Exec(b.String(), vals...)
	if err != nil {
		return err
	}
	return nil
}

func (q *GenreQueries) GetGenreMap() (map[string]int, error) {
	genreMap := map[string]int{}
	query := "SELECT Id, Name FROM library_genre"

	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		genre := models.Genre{}
		if err := rows.Scan(&genre.Id, &genre.Name); err != nil {
			return nil, err
		}
		genreMap[genre.Name] = genre.Id
	}
	return genreMap, nil
}
