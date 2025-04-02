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
	// Periksa status koneksi database
	if config.ErrorMongoconn != nil || config.ErrorMongoconnpaperka != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"message": "Koneksi database gagal"})
	}

	// Ambil profil Athena
	profileAthena, err := mgdb.GetOneDoc[model.Profile](config.Mongoconn, "profile", bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mendapatkan profil Athena", "error": err.Error()})
	}
	stata, resa, erra := RefreshToken(profileAthena, false)

	// Ambil profil Paperka
	profilePaperka, err := mgdb.GetOneDoc[model.Profile](config.Mongoconnpaperka, "profile", bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mendapatkan profil Paperka", "error": err.Error()})
	}
	statb, resb, errb := RefreshToken(profilePaperka, true)

	// Rekap laporan Satpam setiap tanggal 1
	laporanSatpam := satpam.ReportBulanKemarin(profileAthena)

	// Kembalikan respons dalam format JSON
	return c.JSON(fiber.Map{
		"httpAthena":   stata,
		"resAthena":    resa,
		"errAthena":    erra,
		"httpPaperka":  statb,
		"urlAthena":    profileAthena.URL,
		"urlPaperka":   profilePaperka.URL,
		"resPaperka":   resb,
		"errPaperka":   errb,
		"reportSatpam": laporanSatpam,
	})
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
