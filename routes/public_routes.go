package route

import (
	"github.com/gofiber/fiber/v2"
	handler "github.com/quocquann/locallibrary/handlers"
)

func PublicRoute(app *fiber.App) {
	bookHandler := handler.NewBookHandler()
	v1 := app.Group("/api/v1")
	{
		v1.Get("/books", bookHandler.GetBooks)
	}
}
