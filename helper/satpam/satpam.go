package satpam

import (
	"strconv"
	"time"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/mod/presensi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RekapBulanBerjalanPerPhoneNumber(phonenumber string) (masuk, pulang int64) {
	filter := FilterBulanBerjalanDenganPhoneNumber(phonenumber)
	pulang, _ = mgdb.GetCountDoc(config.Mongoconn, "logpresensi", filter)
	masuk, _ = mgdb.GetCountDoc(config.Mongoconn, "selfie", filter)
	return
}

func RekapRatesBulanBerjalan(db *mongo.Database, phonenumber string) (RekapRating, error) {
	filter := FilterBulanBerjalanDenganPhoneNumber(phonenumber)

	docs, err := mgdb.GetAllDoc[[]presensi.PresensiSelfie](db, "selfie", filter)
	if err != nil {
		return RekapRating{}, err
	}

	total := 0
	count := 0
	jumlahPerRating := make(map[int]int)

	for _, doc := range docs {
		for _, rate := range doc.Rates {
			total += rate.Rating
			count++
			jumlahPerRating[rate.Rating]++
		}
	}

	avg := 0.0
	if count > 0 {
		avg = float64(total) / float64(count)
	}

	return RekapRating{
		PhoneNumber:     phonenumber,
		TotalRating:     total,
		AverageRating:   avg,
		JumlahPerRating: jumlahPerRating,
	}, nil
}

func ReportBulanKemarin(profile model.Profile) string {
	// Ambil tanggal hari ini saja dan satu kali saja dijalankan
	if time.Now().Day() != 1 {
		return "bukan tanggal 1"
	}
	countlog, _ := mgdb.GetCountDoc(config.Mongoconn, "logreportbulan", FilterHariIni())
	if countlog > 0 {
		return "sudah dilakukan rekap pada hari ini sebelumnya"
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
		//dt := &model.TextMessage{
		//	To:       satpam.Phonenumber,
		//	IsGroup:  false,
		//	Messages: msg,
		//}
		//go jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	}
	dt := &model.TextMessage{
		To:       profile.WAGroupWarga,
		IsGroup:  true,
		Messages: msg,
	}
	httpcode, res, err := jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	if err != nil {
		return "gagal kirim pesan ke grup warga: " + err.Error()
	}
	var lap LaporanBulanan
	lap.Message = msg
	id, err := mgdb.InsertOneDoc(config.Mongoconn, "logreportbulan", lap)
	if err != nil {
		return err.Error() + " Error ketika insert ke logreportbulan"
	}
	go KirimReportKeSatpam(msg, satpams, profile)
	return "berhasil insert: " + id.Hex() + "|" + strconv.Itoa(httpcode) + "|" + res.Response + "|" + res.Info
}

func KirimReportKeSatpam(msg string, satpams []model.Satpam, profile model.Profile) {
	for _, satpam := range satpams {
		dt := &model.TextMessage{
			To:       satpam.Phonenumber,
			IsGroup:  false,
			Messages: msg,
		}
		jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, config.APIWAText)
	}
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

// selfie dan logpresensi
func FilterBulanBerjalanDenganPhoneNumber(phonenumber string) (filter bson.M) {
	// Ambil waktu sekarang
	currentTime := time.Now()

	// Hitung awal bulan berjalan (bulan ini)
	startOfCurrentMonth := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Akhir periode adalah waktu saat ini
	endOfCurrentMonth := currentTime

	// Konversi ke ObjectID berdasarkan timestamp
	startObjectID := primitive.NewObjectIDFromTimestamp(startOfCurrentMonth)
	endObjectID := primitive.NewObjectIDFromTimestamp(endOfCurrentMonth)

	// Buat filter MongoDB
	filter = bson.M{
		"_id": bson.M{
			"$gte": startObjectID,
			"$lte": endObjectID,
		},
		"phonenumber": phonenumber,
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
