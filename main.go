package main

import (
	"github.com/gocroot/lite/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(Config)
	app.Use(cors.New(Cors))
	route.Web(app)
	app.Listen(IPPort)
}
