package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/quocquann/locallibrary/models"
)

func CrawlBook(url string) ([]models.Book, error) {
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
		image := s.Find(".product-thumbnail>a.image_link.display_flex>img").AttrOr("data-lazyload", "")
		detailUrl := s.Find("h3.product-name a").AttrOr("href", "")
		ch := make(chan string, 2)
		go getDetail(url+detailUrl, ch)
		book := models.Book{Title: title, Image: "https:" + image, Author: <-ch, Genre: <-ch}
		books = append(books, book)
	})

	return books, nil
}

func getDetail(url string, ch chan string) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return
	}

	doc.Find(".group-status").Each(func(i int, s *goquery.Selection) {
		author := s.Find("p:first-child a").Text()
		genre := s.Find("p:nth-child(2) a").Text()
		ch <- author
		ch <- genre
	})
}
