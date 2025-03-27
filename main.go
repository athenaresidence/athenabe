package main

import (
	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/bot"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper"
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
		//refresh token
		profile, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconnpaperka, "profile", bson.M{})
		var wh model.WebHook
		wh.Secret = config.PaperkaSecret
		wh.URL = profile.URL
		wh.ReadStatusOff = true
		stat, userwa, err := jsonapi.PostStructWithToken[model.User]("token", profile.Token, wh, config.APISignUp)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error(), "stat": stat})
		}
		if stat == 200 {
			res, err := mgdb.UpdateOneDoc(config.Mongoconnpaperka, "profile", bson.M{"secret": config.PaperkaSecret}, bson.M{"token": userwa.Token})
			if err != nil {
				return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": err.Error()})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": res})
		} else {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"stat": stat})
		}
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
			go jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	})

	app.Post("/webhook/athena", func(c *fiber.Ctx) error {
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
		waphonenumber := "628176900300"
		prof, err := helper.GetAppProfile(waphonenumber, config.Mongoconn)
		if err != nil {
			resp.Response = err.Error()
			return c.Status(fiber.StatusServiceUnavailable).JSON(resp)
		}

		if msg.Message != "" {
			resp, err = helper.WebHook(prof.QRKeyword, waphonenumber, config.WAAPIQRLogin, config.APIWAText, msg, config.Mongoconn)
			if err != nil {
				resp.Response = err.Error()
			}
			return c.Status(fiber.StatusOK).JSON(resp)
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	})

	app.Listen(IPPort)
}
