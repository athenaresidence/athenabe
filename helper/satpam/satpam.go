package satpam

import (
	"strconv"
	"time"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReportBulanKemarin(profile model.Profile) bool {
	// Ambil tanggal hari ini saja dan satu kali saja dijalankan
	if time.Now().Day() != 1 {
		return false
	}
	countlog, _ := mgdb.GetCountDoc(config.Mongoconn, "logreportbulan", FilterHariIni())
	if countlog > 0 {
		return false
	}

	//ambil data satpam
	satpams, _ := mgdb.GetAllDoc[[]model.Satpam](config.Mongoconn, "satpam", bson.M{})
	msg := "*Rekapitulasi Kehadiran Satpam Bulan Kemarin:*\n"
	for i, satpam := range satpams {
		msg += strconv.Itoa(i+1) + ". " + satpam.Nama + " : "
		// Filter berdasarkan _id dan phonenumber
		filter := FilterBulanKemarendenganPhoneNumber(satpam.Phonenumber)
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
	var lap LaporanBulanan
	lap.Message = msg
	mgdb.InsertOneDoc(config.Mongoconn, "logreportbulan", lap)
	return true
}

func FilterBulanKemarendenganPhoneNumber(phonenumber string) (filter bson.M) {
	// Hitung awal dan akhir bulan kemarin
	today := time.Now()
	lastMonth := today.AddDate(0, -1, 0)
	startOfLastMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfLastMonth := startOfLastMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Konversi ke ObjectID (ID MongoDB pertama dan terakhir bulan lalu)
	startObjectID := primitive.NewObjectIDFromTimestamp(startOfLastMonth)
	endObjectID := primitive.NewObjectIDFromTimestamp(endOfLastMonth)
	filter = bson.M{
		"_id": bson.M{
			"$gte": startObjectID,
			"$lte": endObjectID,
		},
		"phonenumber": phonenumber, // Ganti dengan nomor yang ingin difilter
	}
	return

}

func FilterHariIni() bson.M {
	// Hitung awal dan akhir hari ini
	today := time.Now()
	startOfToday := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfToday := startOfToday.Add(24*time.Hour - time.Nanosecond)

	// Konversi ke ObjectID berdasarkan timestamp
	startObjectID := primitive.NewObjectIDFromTimestamp(startOfToday)
	endObjectID := primitive.NewObjectIDFromTimestamp(endOfToday)

	// Filter berdasarkan _id (hari ini saja)
	filter := bson.M{
		"_id": bson.M{
			"$gte": startObjectID,
			"$lte": endObjectID,
		},
	}
	return filter
}
