package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/quocquann/locallibrary/crawler"
)

func mockHttpServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestCrawlFunc(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		responseData, err := os.ReadFile("./fixtures/res.txt")
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}

	mockServer := mockHttpServer(mockHandler)
	defer mockServer.Close()

	url := mockServer.URL

	books, err := crawler.CrawlBook(url)

	if err != nil {
		t.Fatalf("error: %v", err)
	}

	const expect int = 20
	if expect != len(books) {
		t.Fatalf("expect %d book(s), got %d book(s)", expect, len(books))
	}
}
