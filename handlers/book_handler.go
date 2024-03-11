package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quocquann/locallibrary/crawler"
	"github.com/quocquann/locallibrary/db"
	"github.com/quocquann/locallibrary/models"
	"github.com/quocquann/locallibrary/queries"
)

type IBookHandler interface {
	// GetBooks(c *fiber.Ctx) error
	GetBooks(c *fiber.Ctx) error
}

func NewBookHandler() IBookHandler {
	return &BookHandler{}
}

type BookHandler struct{}

func (*BookHandler) GetBooks(c *fiber.Ctx) error {
	db, err := db.OpenConnection()
	if err != nil {
		return err
	}
	bookQueries := queries.NewBookQueries(db)
	titleQuery := c.Query("title", "")
	authorQuery := c.Query("author", "")
	lastTimeInsertedBookString, err := bookQueries.GetLastTimeInsertedBook()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"data":  nil,
		})
	}

	lastTimeInsertedBook, err := time.Parse("2006-01-02 15:04:05", lastTimeInsertedBookString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"data":  nil,
		})
	}

	if time.Since(lastTimeInsertedBook) < time.Hour*24 {
		fmt.Println("FROM DATABASE")
		books, err := bookQueries.GetBooks(titleQuery, authorQuery)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
				"data":  nil,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"msg":   "success",
			"data":  books,
		})
	}

	fmt.Println("FROM Crawl")
	books, err := crawler.CrawlBook("https://gacxepbookstore.vn")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"data":  nil,
		})
	}

	if titleQuery != "" {
		b := []models.Book{}
		for _, book := range books {
			if strings.Contains(strings.ToLower(book.Title), titleQuery) {
				b = append(b, book)
			}
		}
		books = b
	}

	if err := bookQueries.AddBooks(books); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"data":  nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  books,
	})
}
