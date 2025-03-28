package controller

import (
	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Homepage(c *fiber.Ctx) error {
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
	}
	//refresh token athena
	profile, _ = mgdb.GetOneDoc[model.Profile](config.Mongoconn, "profile", bson.M{})
	wh.Secret = config.PaperkaSecret
	wh.URL = profile.URL
	wh.ReadStatusOff = false
	//wh.SendTyping = true
	stat, userwa, err = jsonapi.PostStructWithToken[model.User]("token", profile.Token, wh, config.APISignUp)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error(), "stat": stat})
	}
	if stat == 200 {
		res, err := mgdb.UpdateOneDoc(config.Mongoconn, "profile", bson.M{"secret": config.PaperkaSecret}, bson.M{"token": userwa.Token})
		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": res})
	} else {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"stat": stat})
	}

}
