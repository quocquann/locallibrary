package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	//TODO: this line use for home page
	// bookItems := doc.Find(".product-loop-1.product-loop-sea.product-base")
	//TODO: this line use for all book page
	bookItems := doc.Find(".product-loop-1.product-base")
	numBookItem := bookItems.Length()

	const numWorker int = 20
	jobs := make(chan models.BookBaseInfo, numBookItem)
	result := make(chan models.Book)

	bookItems.Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3.product-name a").Text()
		// image := s.Find(".product-thumbnail>a.image_link.display_flex>img").AttrOr("data-lazyload", "")
		image := s.Find("h3.product-name a").AttrOr("href", "")

		detailUrl := s.Find("h3.product-name a").AttrOr("href", "")
		//TODO: this line use for all book page]
		url = strings.Split(url, "/all")[0]
		jobs <- models.BookBaseInfo{Title: title, Image: image, Url: url + detailUrl}
	})

	close(jobs)

	for i := 0; i < numWorker; i++ {
		go getDetail(jobs, result)
	}

	for i := 0; i < numBookItem; i++ {
		res := <-result
		if res.Isbn != "" && res.Publisher.Name != "" {
			books = append(books, res)
		}
	}

	return books, nil
}

func getDetail(jobs chan models.BookBaseInfo, result chan models.Book) {
	for job := range jobs {
		fmt.Println(job.Url)
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
		isbn := doc.Find(".rte.description.rte-summary ul li:nth-child(4)").Text()
		publisher := doc.Find(".rte.description.rte-summary ul li:nth-child(1)").Text()
		describe := doc.Find("#tabs-1 p:first-child").Text()

		if !strings.HasPrefix(isbn, "ISBN") {
			isbn = ""
		} else {
			isbn = isbn[len(isbn)-10:]
		}

		if !strings.HasPrefix(publisher, "Publisher") {
			publisher = ""
		} else {
			fmt.Println()
			publisher = strings.TrimSpace(strings.Split(publisher, ":")[1])
		}
		result <- models.Book{Isbn: isbn, Title: job.Title, Image: job.Image, Describe: describe, Author: models.Author{Name: author}, Genre: models.Genre{Name: genre}, Publisher: models.Publisher{Name: publisher}}
	}
}
