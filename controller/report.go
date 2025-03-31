package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReportBulananSatpam(c *fiber.Ctx) error {
	// Ambil tanggal hari ini
	today := time.Now()
	if today.Day() != 1 {
		return c.Status(fiber.StatusOK).JSON(today)
	}
	// Hitung awal dan akhir bulan kemarin
	lastMonth := today.AddDate(0, -1, 0)
	startOfLastMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfLastMonth := startOfLastMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Konversi ke ObjectID (ID MongoDB pertama dan terakhir bulan lalu)
	startObjectID := primitive.NewObjectIDFromTimestamp(startOfLastMonth)
	endObjectID := primitive.NewObjectIDFromTimestamp(endOfLastMonth)

	//ambil data satpam
	satpams, _ := mgdb.GetAllDoc[[]model.Satpam](config.Mongoconn, "satpam", bson.M{})
	msg := "*Rekapitulasi Kehadiran Satpam Bulan Kemarin:*\n"
	for i, satpam := range satpams {
		msg += strconv.Itoa(i+1) + ". " + satpam.Nama + " : "
		fmt.Println("Nilai:", satpam)
		// Filter berdasarkan _id dan phonenumber
		filter := bson.M{
			"_id": bson.M{
				"$gte": startObjectID,
				"$lte": endObjectID,
			},
			"phonenumber": satpam.Phonenumber, // Ganti dengan nomor yang ingin difilter
		}
		count, _ := mgdb.GetCountDoc(config.Mongoconn, "logpresensi", filter)
		msg += strconv.Itoa(int(count)) + "\n"
	}

	return c.Status(fiber.StatusOK).JSON(satpams)
}
