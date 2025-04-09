package presensi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CekSelfiePulang(Profile model.Profile, Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.. " + Pesan.Alias_name
	}

	filter := bson.M{"_id": mgdb.TodayFilter(), "phonenumber": Pesan.Phone_number} //, "ismasuk": false}
	pstoday, err := mgdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf kakak belum share live location dulu, silahkan share live location / berbagi lokasi terkini dengan ditambah caption\n*pulang*"
	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     false,
	}
	selfistat, err := mgdb.InsertOneDoc(db, "selfie", pselfie)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan input ke database " + err.Error()
	}
	filter = bson.M{"_id": mgdb.TodayFilter(), "cekinlokasi.phonenumber": Pesan.Phone_number, "ismasuk": true}
	selfiemasuk, err := mgdb.GetOneLatestDoc[PresensiSelfie](db, "selfie", filter)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf kakak belum selfie masuk. " + err.Error()
	}
	// Ekstrak timestamp dari ObjectID
	objectIDTimestamp := selfiemasuk.ID.Timestamp()
	// Dapatkan waktu saat ini
	currentTime := time.Now()
	// Hitung selisih waktu dalam detik
	diff := currentTime.Sub(objectIDTimestamp) //.Seconds()
	// Konversi selisih waktu ke jam, menit, dan detik
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60
	KetJam := fmt.Sprintf("%d jam, %d menit, %d detik", hours, minutes, seconds)

	skor := (diff.Seconds() / 43200) * 100 //selisih waktu dibagi 12 jam
	skorValue := fmt.Sprintf("%f", skor)
	//post ke backedn domyikado
	datapresensi := PresensiDomyikado{
		ID:          selfiemasuk.ID,
		PhoneNumber: Pesan.Phone_number,
		Skor:        skor,
		KetJam:      KetJam,
		LamaDetik:   diff.Seconds(),
		Lokasi:      pstoday.Lokasi.Nama,
	}
	//kalo satpam maka kirim ke grup dan simpan database
	satpam, err := mgdb.GetOneDoc[model.Satpam](config.Mongoconn, "satpam", bson.M{"phonenumber": Pesan.Phone_number})
	if err != mongo.ErrNoDocuments {
		go FaceDetectLeafly(Profile, Pesan, db, selfistat) //kirim deteksi leafly
		msg := "*Pulang Shift Jaga*\n" + satpam.Nama + "\n" + satpam.Phonenumber + "\nHadir selama: " + KetJam + "\n*Skor: " + skorValue + "*"
		notifgroup := model.ImageMessage{
			To:          Profile.WAGroupWarga,
			IsGroup:     true,
			Base64Image: Pesan.Filedata,
			Caption:     msg,
		}
		datapresensi.Nama = satpam.Nama
		stat, resp, err := jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifgroup, config.APIWAIMG)
		if stat != 200 {
			return "Ada kesalahan pengiriman notif ke grup\n" + notifgroup.To + "\n" + notifgroup.Caption + "\n" + notifgroup.Base64Image + "\n" + err.Error() + "\n" + resp.Response
		}
		mgdb.InsertOneDoc(config.Mongoconn, "logpresensi", datapresensi)
	}
	return "Hai kak, " + Pesan.Alias_name + "\nBerhasil Presensi Pulang di lokasi:" + pstoday.Lokasi.Nama + "\nHadir selama: " + KetJam + "\n*Skor: " + skorValue + "*"

}

func CekSelfieMasuk(Profile model.Profile, Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.. " + Pesan.Alias_name
	}
	filter := bson.M{"_id": mgdb.TodayFilter(), "phonenumber": Pesan.Phone_number, "ismasuk": true}
	pstoday, err := mgdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf kakak belum share live location dulu, silahkan share live location / berbagi lokasi terkini dengan ditambah caption\n*masuk*"
	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     true,
	}
	//kalo satpam maka kirim ke grup dan diri sendiri
	satpam, err := mgdb.GetOneDoc[model.Satpam](config.Mongoconn, "satpam", bson.M{"phonenumber": Pesan.Phone_number})
	if err != mongo.ErrNoDocuments {
		pselfie.Nama = satpam.Nama
		pselfie.PhoneNumber = satpam.Phonenumber
		selfistat, err := mgdb.InsertOneDoc(db, "selfie", pselfie)
		if err != nil {
			return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan input ke database " + err.Error()
		}
		go FaceDetectLeafly(Profile, Pesan, db, selfistat) //kirim deteksi leafly
		//ambil rekap bulan berjalan
		jmlmasuk, jmlpulang := RekapBulanBerjalanPerPhoneNumber(satpam.Phonenumber)
		rekaprating, err := RekapRatesBulanBerjalan(config.Mongoconn, satpam.Phonenumber)
		var ratemsg string
		if err == nil {
			ratemsg = fmt.Sprintf(
				"*Feedback _%s_:*\nTotal Rating: %d\nRata-rata: %.2f\nDetail:\n%v",
				satpam.Nama,
				rekaprating.TotalRating,
				rekaprating.AverageRating,
				FormatJumlahPerRating(rekaprating.JumlahPerRating),
			)
		}
		msg := "*Masuk Shift Jaga Satpam Sekarang*\nNama: " + satpam.Nama +
			"\nTelepon: " + satpam.Phonenumber +
			"\n" + ratemsg +
			"\nRekap Presensi Bulan ini: " +
			"\nMasuk: " + strconv.FormatInt(jmlmasuk, 10) +
			"\nPulang: " + strconv.FormatInt(jmlpulang, 10) +
			"\n\nMohon berikan feedback pekerjaan selama shift jaga berjalan ke:\nhttps://athenaresidence.github.io/rate/#" + selfistat.Hex() +
			"\n\n> Mari jadikan komplek kita lebih?"
		notifgroup := model.ImageMessage{
			To:          Profile.WAGroupWarga,
			IsGroup:     true,
			Base64Image: Pesan.Filedata,
			Caption:     msg,
		}
		stat, resp, err := jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifgroup, config.APIWAIMG)
		if stat != 200 {
			return "Ada kesalahan pengiriman notif ke grup\ngrup: " + notifgroup.To + "\ncaption: " + notifgroup.Caption + "\n" + err.Error() + "\n" + resp.Response
		}
		return "Hai kak, " + Pesan.Alias_name +
			"\nCekin Masuk di lokasi: " + pstoday.Lokasi.Nama +
			"\n" + ratemsg +
			"\nRekap Presensi Bulan ini: " +
			"\nMasuk: " + strconv.FormatInt(jmlmasuk, 10) +
			"\nPulang: " + strconv.FormatInt(jmlpulang, 10) +
			"\n> *Jangan lupa share live loc dengan caption *pulang* setelah selesai shift ya kak, biar dianggap masuk shift jaga*"
	}
	return "Hai kak, " + Pesan.Alias_name + ". Mohon maaf nomor kakak belum terdaftar di sistem kita, hubungin admin ya kak.\nCekin Masuk di lokasi: " + pstoday.Lokasi.Nama
}

func FaceDetectLeafly(Profile model.Profile, Pesan model.WAMessage, db *mongo.Database, idselfie primitive.ObjectID) string {
	conf, err := mgdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": Profile.Phonenumber})
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	//sending foto ke leafly, ambil keterangan analisisnya
	dt := FaceDetect{
		IDUser:    Pesan.Phone_number,
		Base64Str: Pesan.Filedata,
	}
	var leaflyfaceket string
	statuscode, faceinfo, err := jsonapi.PostStructWithToken[FaceInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan pemanggilan API leafly :" + err.Error() + " Mohon kirim ulang foto kembali dengan caption *masuk*"
	}
	if statuscode != http.StatusOK {
		if statuscode == http.StatusFailedDependency { //deteksi penggunaan foto yang sama
			leaflyfaceket = "Ubur-ubur ikan lele\n_foto kemarin le_"
		} else if statuscode == http.StatusMultipleChoices { //deteksi wajah lebih dari 1 : status 300
			leaflyfaceket = "Ubur-ubur ikan lele\n_ada penampakan le_"
			return ""
		} else { //tidak terdeteksi wajah : 410 status
			leaflyfaceket = "Ubur-ubur ikan lele\n_Wajah rata bre_"
			if faceinfo.FileHash == "" {
				faceinfo.FileHash = dt.Base64Str
			}
			//return "Wah kak " + Pesan.Alias_name + " mohon maaf:\n" + faceinfo.Error + "\nCode: " + strconv.Itoa(statuscode)
		}
	}
	go mgdb.UpdateOneDoc(db, "selfie", bson.M{"_id": idselfie}, bson.M{"iduser": faceinfo.PhoneNumber,
		"commit":     faceinfo.Commit,
		"filehash":   faceinfo.FileHash,
		"facestatus": statuscode,
		"remaining":  faceinfo.Remaining,
	})
	pesandeteksi := model.TextMessage{
		To:       Profile.WAGroupWarga,
		IsGroup:  true,
		Messages: leaflyfaceket,
	}
	jsonapi.PostStructWithToken[model.Response]("Token", Profile.Token, pesandeteksi, config.APIWAText)
	return leaflyfaceket
}

func PresensiMasuk(Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if !Pesan.LiveLoc {
		return "Minimal share live location dulu lah kak " + Pesan.Alias_name
	}
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak, kakak " + Pesan.Alias_name + " belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru berbagi lokasi terkini / share live location dengan caption *masuk*"
	}
	if lokasiuser.Nama == "" {
		return "Nama nya kosong kak " + Pesan.Alias_name
	}
	dtuser := &PresensiLokasi{
		PhoneNumber: Pesan.Phone_number,
		Lokasi:      lokasiuser,
		IsMasuk:     true,
		CreatedAt:   time.Now(),
	}
	_, err = mgdb.InsertOneDoc(db, "presensi", dtuser)
	if err != nil {
		return "Gagal insert ke database kak " + Pesan.Alias_name
	}

	return "Hai.. hai.. kakak atas nama:\n*" + Pesan.Alias_name + "*\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nsilahkan dilanjutkan dengan selfie di lokasi ya maximal 5 menit setelah share live location, jangan lupa ditambah caption\n*masuk*"
}

func PresensiPulang(Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if !Pesan.LiveLoc {
		return "Minimal share live location dulu lah kak " + Pesan.Alias_name
	}
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak " + Pesan.Alias_name + ", kakak belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru cekin pulang."
	}
	if lokasiuser.Nama == "" {
		return "Nama nya kosong kak " + Pesan.Alias_name
	}
	dtuser := &PresensiLokasi{
		PhoneNumber: Pesan.Phone_number,
		Lokasi:      lokasiuser,
		IsMasuk:     false,
		CreatedAt:   time.Now(),
	}
	filter := bson.M{"_id": mgdb.TodayFilter(), "cekinlokasi.phonenumber": Pesan.Phone_number, "ismasuk": true}
	docselfie, err := mgdb.GetOneLatestDoc[PresensiSelfie](db, "selfie", filter)
	if err != nil {
		return "Kakak " + Pesan.Alias_name + " belum selfie masuk ini " + err.Error()
	}
	if docselfie.CekInLokasi.Lokasi.ID != lokasiuser.ID {
		return "Lokasi pulang nya harus sama dengan lokasi masuknya kak " + Pesan.Alias_name + ".\nLokasi : " + lokasiuser.Nama
	}
	_, err = mgdb.InsertOneDoc(db, "presensi", dtuser)
	if err != nil {
		return "Gagal insert ke database kak " + Pesan.Alias_name
	}

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nsilahkan dilanjutkan dengan selfie di lokasi ya maximal 5 menit setelah share live location, jangan lupa ditambah caption\n*pulang*"
}

func GetLokasi(mongoconn *mongo.Database, long float64, lat float64) (lokasi Lokasi, err error) {
	filter := bson.M{
		"batas": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}

	lokasi, err = mgdb.GetOneDoc[Lokasi](mongoconn, "lokasi", filter)
	if err != nil {
		return
	}
	return
}
