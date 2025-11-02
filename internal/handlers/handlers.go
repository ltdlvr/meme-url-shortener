package handlers

import (
	"database/sql"
	"meme-url/internal/models"
	"meme-url/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "pong",
	})
}

func Shorten(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Request

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid JSON",
			})
		}
		if req.Url == "" { //тут надо будет сделать валидацию
			return c.Status(400).JSON(fiber.Map{
				"error": "URL can't be empty",
			})
		}

		res, err := repository.GetShortLink(db, req.Url)
		if err != nil {
			log.Errorf("failed to create short link: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Couldn't get a short link",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(res)

	}
}

func Redirect(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortcode := c.Params("shortcode")
		if shortcode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "shortcode is required",
			})
		}

		longUrl, err := repository.GetLongURL(db, shortcode)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error": "shortcode not found",
			})
		}
		if err != nil {
			log.Errorf("failed to get long URL: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Redirect(longUrl, fiber.StatusMovedPermanently)
	}
}

func DeleteShortcode(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortcode := c.Params("shortcode")
		if shortcode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "shortcode required",
			})
		}

		err := repository.DeleteShortcode(db, shortcode)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error": "shortcode not found",
			})
		}
		if err != nil {
			log.Errorf("failed to delete shortcode %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "deleted succesfully",
		})
	}
}
