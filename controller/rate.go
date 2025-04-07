package controller

import (
	"strconv"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper"
	"github.com/gocroot/lite/mod/presensi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCekInSelfieData(c *fiber.Ctx) error {
	idselfie := c.Params("idselfie")
	objectId, err := primitive.ObjectIDFromHex(idselfie)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     objectId.String(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	hasil, err := mgdb.GetOneLatestDoc[presensi.PresensiSelfie](config.Mongoconn, "selfie", bson.M{"_id": objectId})
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     objectId.String(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	return c.Status(fiber.StatusOK).JSON(hasil)
}

func PostRateSelfie(c *fiber.Ctx) error {
	var rating presensi.Rating
	err := c.BodyParser(&rating)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     "Body tidak valid",
		}
		return c.Status(fiber.StatusFailedDependency).JSON(respn)
	}
	objectId, err := primitive.ObjectIDFromHex(rating.ID)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     objectId.String(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	hasil, err := mgdb.GetOneLatestDoc[presensi.PresensiSelfie](config.Mongoconn, "selfie", bson.M{"_id": objectId})
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     objectId.String(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	res, err := mgdb.AddDocToArray(config.Mongoconn, "selfie", hasil.ID, "rates", rating)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     "Tidak berhasil update rating ke selfie",
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	//kirim pesan ke yang piket
	prof, err := helper.GetAppProfile(config.AthenaBotNumber, config.Mongoconn)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     "Tidak berhasil mendapatkan profile aplikasi",
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(respn)
	}
	dt := model.TextMessage{
		To:       hasil.PhoneNumber,
		IsGroup:  false,
		Messages: "Anda mendapatkan feedback pekerjaan dari " + rating.Nomor + " dengan rating bintang " + strconv.Itoa(rating.Rating) + "\nDengan feedback:\n" + rating.Komentar,
	}
	go jsonapi.PostStructWithToken[model.Response]("Token", prof.Token, dt, config.APIWAText)
	respn := model.Response{
		Info: strconv.FormatInt(res.ModifiedCount, 10),
	}
	return c.Status(fiber.StatusOK).JSON(respn)
}
