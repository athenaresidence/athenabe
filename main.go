package main

import (
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(Config)
	app.Use(cors.New(Cors))

	app.Get("/", func(c *fiber.Ctx) error {
		//checkstatus koneksi db
		if config.ErrorMongoconn != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": "Koneksi database gagal"})
		}
		if config.ErrorMongoconnpaperka != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": "Koneksi database gagal"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success deploy ulang"})
	})
	app.Post("/webhook", func(c *fiber.Ctx) error {
		var h model.Header
		c.ReqHeaderParser(&h)
		return c.Status(fiber.StatusOK).JSON(h)
	})

	app.Listen(IPPort)
}
