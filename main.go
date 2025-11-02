package main

import (
	"meme-url/internal/db"
	"meme-url/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db := db.InitDatabase()
	defer db.Close()

	app := fiber.New()

	app.Use(logger.New())

	app.Get("/ping", handlers.HealthCheck)
	app.Post("/api/shorten", handlers.Shorten(db))
	app.Get("/:shortcode", handlers.Redirect(db))
	app.Delete("/:shortcode", handlers.DeleteShortcode(db))
	app.Listen(":3030")

}
