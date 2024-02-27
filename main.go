package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/quocquann/locallibrary/models"
)

func crawlBook(url string) ([]models.Book, error) {
	books := []models.Book{}
	res, err := http.Get(url)
	if err != nil {
		return []models.Book{}, err
	}

	if res.StatusCode != 200 {
		return []models.Book{}, fmt.Errorf("status code error: %d", res.StatusCode)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return []models.Book{}, err
	}

	doc.Find(".product-loop-1.product-loop-sea.product-base").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3.product-name a").Text()
		book := models.Book{Title: title}
		books = append(books, book)
	})

	return books, nil
}

func main() {
	books, err := crawlBook("https://gacxepbookstore.vn/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(books)
}
