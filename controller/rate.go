package controller

import (
	"github.com/gocroot/lite/config"
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
	return c.Status(fiber.StatusBadRequest).JSON(hasil)
}

func PostRateSelfie(c *fiber.Ctx) error {
	var rating presensi.Rating
	err := c.BodyParser(&rating)
	if err != nil {
		respn := model.Response{
			Response: err.Error(),
			Info:     "Bady tidak valid",
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
			Info:     rating.ID,
		}
		return c.Status(fiber.StatusBadRequest).JSON(respn)
	}
	return c.Status(fiber.StatusBadRequest).JSON(res)
}
