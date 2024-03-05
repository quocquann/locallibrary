package crawler

import (
	"fmt"
	"log"
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

	bookItems := doc.Find(".product-loop-1.product-loop-sea.product-base")

	numBookItem := bookItems.Length()

	const numWorker int = 20
	jobs := make(chan models.BookBaseInfo, numBookItem)
	result := make(chan models.Book)

	bookItems.Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3.product-name a").Text()
		image := s.Find(".product-thumbnail>a.image_link.display_flex>img").AttrOr("data-lazyload", "")
		detailUrl := s.Find("h3.product-name a").AttrOr("href", "")
		jobs <- models.BookBaseInfo{Title: title, Image: image, Url: url + detailUrl}
	})

	close(jobs)

	for i := 0; i < numWorker; i++ {
		go getDetail(jobs, result)
	}

	for i := 0; i < numBookItem; i++ {
		books = append(books, <-result)
	}

	return books, nil
}

func getDetail(jobs chan models.BookBaseInfo, result chan models.Book) {
	for job := range jobs {
		res, err := http.Get(job.Url)
		if err != nil {
			log.Println(err)
			return
		}

		if res.StatusCode != 200 {
			log.Printf("status code error: %d\n", res.StatusCode)
			return
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			log.Println(err)
			return
		}

		author := doc.Find(".group-status").First().Find("p:first-child a").Text()
		genre := doc.Find(".group-status").First().Find("p:nth-child(2) a").Text()
		result <- models.Book{Title: job.Title, Image: "https:" + job.Image, Author: author, Genre: genre}
	}
}
