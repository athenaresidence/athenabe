package main

import (
	"github.com/gocroot/lite/bot"
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
		err := c.ReqHeaderParser(&h)
		if err != nil {
			return c.Status(fiber.StatusFailedDependency).JSON(err)
		}
		resp := model.Response{Response: h.Secret}
		if h.Secret != config.PaperkaSecret {
			return c.Status(fiber.StatusForbidden).JSON(resp)
		}
		var msg model.WAMessage
		err = c.BodyParser(&msg)
		if err != nil {
			return c.Status(fiber.StatusFailedDependency).JSON(err)
		}
		if !msg.Is_group {
			profile, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconnpaperka, "profile", bson.M{})
			dt := &model.TextMessage{
				To:      msg.Chat_number,
				IsGroup: msg.Is_group,
			}
			dt.Messages = bot.HandlerPesan(msg, profile)
			atapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	})

	app.Listen(IPPort)
}
