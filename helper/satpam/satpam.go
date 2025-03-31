package satpam

import (
	"strconv"
	"time"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReportBulanKemarin() bool {
	// Ambil tanggal hari ini
	today := time.Now()
	if today.Day() != 1 {
		return false
	}
	profile, _ := helper.GetAppProfile(config.AthenaBotNumber, config.Mongoconn)
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
		// Filter berdasarkan _id dan phonenumber
		filter := bson.M{
			"_id": bson.M{
				"$gte": startObjectID,
				"$lte": endObjectID,
			},
			"phonenumber": satpam.Phonenumber, // Ganti dengan nomor yang ingin difilter
		}
		count, _ := mgdb.GetCountDoc(config.Mongoconn, "logpresensi", filter)
		msg += strconv.Itoa(int(count)) + " shift\n"
		dt := &model.TextMessage{
			To:       satpam.Phonenumber,
			IsGroup:  false,
			Messages: msg,
		}
		go jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	}
	dt := &model.TextMessage{
		To:       profile.WAGroupWarga,
		IsGroup:  true,
		Messages: msg,
	}
	go jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	return true
}
