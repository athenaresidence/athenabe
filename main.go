package main

import (
	"github.com/gocroot/lite/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(Config)
	app.Use(cors.New(Cors))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})
	app.Post("/webhook", func(c *fiber.Ctx) error {
		var h model.Header
		c.ReqHeaderParser(&h)
		return c.Status(fiber.StatusOK).JSON(h)
	})

	app.Listen(IPPort)
}
