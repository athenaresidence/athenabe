package controller

import (
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper"
	"github.com/gocroot/lite/model"
	"github.com/gofiber/fiber/v2"
)

func BotAthena(c *fiber.Ctx) error {
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
	waphonenumber := c.Params("waphonenumber")
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

}
