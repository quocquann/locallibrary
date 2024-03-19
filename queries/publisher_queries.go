package queries

import (
	"bytes"
	"database/sql"

	"github.com/quocquann/locallibrary/models"
)

type IPublisherQueries interface {
	AddPublisher([]models.Book) error
	GetPublisherMap() (map[string]int, error)
}

func NewPublisherQueries(db *sql.DB) IPublisherQueries {
	return &PublisherQueries{
		DB: db,
	}
}

type PublisherQueries struct {
	*sql.DB
}

func (q *PublisherQueries) AddPublisher(books []models.Book) error {
	if len(books) <= 0 {
		return nil
	}

	query := "INSERT INTO library_publisher(Name) VALUES "
	b := bytes.Buffer{}
	b.WriteString(query)
	vals := []interface{}{}

	for _, book := range books {
		b.WriteString("(?),")
		vals = append(vals, book.Publisher.Name)
	}

	b.Truncate(b.Len() - 1)
	b.WriteString(" ON DUPLICATE KEY UPDATE Name=VALUES(Name)")

	_, err := q.Exec(b.String(), vals...)
	if err != nil {
		return err
	}

	return nil
}

func (q *PublisherQueries) GetPublisherMap() (map[string]int, error) {
	publisherMap := map[string]int{}

	query := "SELECT Id, Name FROM library_publisher"
	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		publisher := models.Publisher{}
		if err := rows.Scan(&publisher.Id, &publisher.Name); err != nil {
			return nil, err
		}
		publisherMap[publisher.Name] = publisher.Id
	}

	return publisherMap, nil
}
