package controller

import (
	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/bot"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func WebHookPaperka(c *fiber.Ctx) error {
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
		profile, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconn, "profile", bson.M{})
		dt := &model.TextMessage{
			To:      msg.Chat_number,
			IsGroup: msg.Is_group,
		}
		dt.Messages = bot.HandlerPesan(msg, profile)
		go jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
