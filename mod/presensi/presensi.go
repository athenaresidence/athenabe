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
	"go.mongodb.org/mongo-driver/mongo"
)

func CekSelfiePulang(Profile model.Profile, Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.. " + Pesan.Alias_name
	}
	dt := FaceDetect{
		IDUser:    Pesan.Phone_number,
		Base64Str: Pesan.Filedata,
	}
	filter := bson.M{"_id": mgdb.TodayFilter(), "phonenumber": Pesan.Phone_number} //, "ismasuk": false}
	pstoday, err := mgdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf kakak belum cekin share live location hari ini " + err.Error()
	}
	conf, err := mgdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": Profile.Phonenumber})
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	statuscode, faceinfo, err := jsonapi.PostStructWithToken[FaceInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan pemanggilan API leafly " + err.Error()
	}
	if statuscode != http.StatusOK {
		if statuscode == http.StatusFailedDependency {
			return "Wah kak " + Pesan.Alias_name + " mohon maaf, jangan kaku gitu dong. Ekspresi wajahnya ga boleh sama dengan selfie sebelumnya ya kak. Senyumnya yang lebar, giginya dilihatin, matanya pelototin, hidungnya keatasin.\n\n" + faceinfo.Error
		} else if statuscode == http.StatusMultipleChoices {
			dt := &model.ImageMessage{
				To:          Pesan.Chat_number,
				Base64Image: faceinfo.FileHash,
				Caption:     faceinfo.Error,
				IsGroup:     Pesan.Is_group,
			}
			statuscode, httpresp, err := jsonapi.PostStructWithToken[model.Response]("Token", Profile.Token, dt, Profile.URLAPIImage)
			if err != nil {
				strconv.Itoa(statuscode)
				return "Akses ke endpoint whatsaut gagal: " + err.Error() + strconv.Itoa(statuscode) + httpresp.Info + httpresp.Response
			}
			return ""
		} else {
			return "Wah kak " + Pesan.Alias_name + " mohon maaf:\n" + faceinfo.Error + "\nCode: " + strconv.Itoa(statuscode)
		}

	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     false,
		IDUser:      faceinfo.PhoneNumber,
		Commit:      faceinfo.Commit,
		Filehash:    faceinfo.FileHash,
		Remaining:   faceinfo.Remaining,
	}
	_, err = mgdb.InsertOneDoc(db, "selfie", pselfie)
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
			return "Ada kesalahan pengiriman notif ke grup\n" + err.Error() + "\n" + resp.Response
		}
		mgdb.InsertOneDoc(config.Mongoconn, "logpresensi", datapresensi)
	}
	return "Hai kak, " + Pesan.Alias_name + "\nBerhasil Presensi Pulang di lokasi:" + pstoday.Lokasi.Nama + "\nHadir selama: " + KetJam + "\n*Skor: " + skorValue + "*"

}

func CekSelfieMasuk(Profile model.Profile, Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if Pesan.Filedata == "" {
		return "Kirim pap nya dulu dong kak.. " + Pesan.Alias_name
	}
	dt := FaceDetect{
		IDUser:    Pesan.Phone_number,
		Base64Str: Pesan.Filedata,
	}
	filter := bson.M{"_id": mgdb.TodayFilter(), "phonenumber": Pesan.Phone_number, "ismasuk": true}
	pstoday, err := mgdb.GetOneDoc[PresensiLokasi](db, "presensi", filter)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf kakak belum cekin share live location hari ini, silahkan share live loc dengan ditambah keyword\n*cekin presensi masuk*\n_" + err.Error() + "_"
	}
	conf, err := mgdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": Profile.Phonenumber})
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}
	statuscode, faceinfo, err := jsonapi.PostStructWithToken[FaceInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan pemanggilan API leafly :" + err.Error()
	}
	if statuscode != http.StatusOK {
		if statuscode == http.StatusFailedDependency {
			return "Wah kak " + Pesan.Alias_name + " mohon maaf, jangan kaku gitu dong. Ekspresi wajahnya ga boleh sama dengan selfie sebelumnya ya kak. Senyumnya yang lebar, giginya dilihatin, matanya pelototin, hidungnya keatasin.\n\n" + faceinfo.Error
		} else if statuscode == http.StatusMultipleChoices {
			dt := &model.ImageMessage{
				To:          Pesan.Chat_number,
				Base64Image: faceinfo.FileHash,
				Caption:     faceinfo.Error,
				IsGroup:     Pesan.Is_group,
			}
			statuscode, httpresp, err := jsonapi.PostStructWithToken[model.Response]("Token", Profile.Token, dt, Profile.URLAPIImage)
			if err != nil {
				strconv.Itoa(statuscode)
				return "Akses ke endpoint whatsaut gagal: " + err.Error() + strconv.Itoa(statuscode) + httpresp.Info + httpresp.Response
			}
			return ""
		} else {
			return "Wah kak " + Pesan.Alias_name + " mohon maaf:\n" + faceinfo.Error + "\nCode: " + strconv.Itoa(statuscode)
		}

	}
	pselfie := PresensiSelfie{
		CekInLokasi: pstoday,
		IsMasuk:     true,
		IDUser:      faceinfo.PhoneNumber,
		Commit:      faceinfo.Commit,
		Filehash:    faceinfo.FileHash,
		Remaining:   faceinfo.Remaining,
	}
	_, err = mgdb.InsertOneDoc(db, "selfie", pselfie)
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan input ke database " + err.Error()
	}
	//kalo satpam maka kirim ke grup
	satpam, err := mgdb.GetOneDoc[model.Satpam](config.Mongoconn, "satpam", bson.M{"phonenumber": Pesan.Phone_number})
	if err != mongo.ErrNoDocuments {
		msg := "*Masuk Shift Jaga*\n" + satpam.Nama + "\n" + satpam.Phonenumber
		notifgroup := model.ImageMessage{
			To:          Profile.WAGroupWarga,
			IsGroup:     true,
			Base64Image: Pesan.Filedata,
			Caption:     msg,
		}
		stat, resp, err := jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifgroup, config.APIWAIMG)
		if stat != 200 {
			return "Ada kesalahan pengiriman notif ke grup\n" + err.Error() + "\n" + resp.Response
		}

	}

	return "Hai kak, " + Pesan.Alias_name + "\nCekin Masuk di lokasi: " + pstoday.Lokasi.Nama + "\n> *Jangan lupa _cekin presensi pulang_ ya kak biar dapat skor*"

}

func PresensiMasuk(Pesan model.WAMessage, db *mongo.Database) (reply string) {
	if !Pesan.LiveLoc {
		return "Minimal share live location dulu lah kak " + Pesan.Alias_name
	}
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)
	lokasiuser, err := GetLokasi(db, Pesan.Longitude, Pesan.Latitude)
	if err != nil {
		return "Mohon maaf kak, kakak " + Pesan.Alias_name + " belum berada di lokasi presensi, silahkan menuju lokasi presensi dahulu baru *cekin presensi masuk*."
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

	return "Hai.. hai.. kakak atas nama:\n*" + Pesan.Alias_name + "*\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nsilahkan dilanjutkan dengan selfie di lokasi ya maximal 5 menit setelah share live location, jangan lupa ditambah keyword\n*selfie presensi masuk*"
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

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + lokasiuser.Nama + "\nsilahkan dilanjutkan dengan selfie di lokasi ya maximal 5 menit setelah share live location, jangan lupa ditambah keyword\n*selfie presensi pulang*"
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
