package controller

import (
	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper/satpam"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Homepage(c *fiber.Ctx) error {
	//checkstatus koneksi db
	if config.ErrorMongoconn != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": "Koneksi database gagal"})
	}
	if config.ErrorMongoconnpaperka != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": "Koneksi database gagal"})
	}
	//dapetin profile apps athena
	profileathena, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconn, "profile", bson.M{})
	stata, resa, erra := RefreshToken(profileathena, false)
	//profile paperka
	profile, _ := mgdb.GetOneDoc[model.Profile](config.Mongoconnpaperka, "profile", bson.M{})
	statb, resb, errb := RefreshToken(profile, true)
	//rekap satpam setiap tanggal 1
	lapket := satpam.ReportBulanKemarin(profileathena)
	//return semuanya
	return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
		"httpathena":   stata,
		"resathena":    resa,
		"errathena":    erra,
		"httppaperka":  statb,
		"respaperka":   resb,
		"errpaperka":   errb,
		"reportsatpam": lapket})

}

func RefreshToken(profile model.Profile, readstatus bool) (stat int, res *mongo.UpdateResult, err error) {
	var wh model.WebHook
	wh.Secret = config.PaperkaSecret
	wh.URL = profile.URL
	wh.ReadStatusOff = readstatus
	stat, userwa, err := jsonapi.PostStructWithToken[model.User]("token", profile.Token, wh, config.APISignUp)
	if err != nil {
		return
	}
	if stat == 200 {
		res, err = mgdb.UpdateOneDoc(config.Mongoconnpaperka, "profile", bson.M{"secret": config.PaperkaSecret}, bson.M{"token": userwa.Token})
	}
	return
}
