package main

import (
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper/atapi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
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
	app.Post("/webhook/paperka", func(c *fiber.Ctx) error {
		var h model.Header
		c.ReqHeaderParser(&h)
		resp := model.Response{Response: h.Secret}
		if h.Secret != config.PaperkaSecret {
			return c.Status(fiber.StatusForbidden).JSON(resp)
		}
		var msg model.WAMessage
		c.BodyParser(&msg)
		if msg.Phone_number == "6281312000300" {
			profile, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconnpaperka, "profile", bson.M{})
			dt := &model.TextMessage{
				To:       msg.Chat_number,
				IsGroup:  msg.Is_group,
				Messages: "ngops... ngops\nkuyy..",
			}
			atapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	})

	app.Listen(IPPort)
}
