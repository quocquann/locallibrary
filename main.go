package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	route "github.com/quocquann/locallibrary/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Fail to load .env file")
	}
	app := fiber.New()
	route.PublicRoute(app)
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
