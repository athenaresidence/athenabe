package bukutamu

import (
	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
)

func BukuTamu(Profile model.Profile, Pesan model.WAMessage) (reply string) {
	tamu := ParsePesanFlexible(Pesan.Message)
	tamu.PhoneNumber = Pesan.Phone_number
	tamu.Nama = Pesan.Alias_name
	//mengambil data satpam yang masuk shift aktif terakhir
	satpam, err := mgdb.GetOneLatestDoc[model.Satpam](config.Mongoconn, "selfie", bson.M{})
	if err != nil {
		return "Mohon maaf saat ini belum ada satpam yang menjaga pos"
	}
	//mengambil data penghuni yang dituju
	penghuniyangdikunjungi, err := mgdb.GetOneDoc[model.Penghuni](config.Mongoconn, "penghuni", bson.M{"nomorrumah": tamu.BlokRumah})
	if err != nil && tamu.Tujuan != "" {
		return "Mohon maaf tujuan yang kakak input tidak valid:\n_" + tamu.Tujuan + "_"
	}
	//pemberitahuan ke grup untuk ada yang jualan masuk ke komplek
	if tamu.Tujuan == "" {
		notifgroup := model.TextMessage{
			To:      Profile.WAGroupWarga,
			IsGroup: true,
			Messages: "Masuk ke lingkungan komplek, Mitra Penjual:\n" + tamu.Nama +
				"\n" + tamu.PhoneNumber +
				"\n" + tamu.Kategori,
		}
		go jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifgroup, config.APIWAText)
		//pemberitahuan ke satpam shift jaga
		notifsatpam := model.TextMessage{
			To:      satpam.Phonenumber,
			IsGroup: false,
			Messages: "Masuk ke lingkungan komplek, Mitra Penjual:\n" + tamu.Nama +
				"\n" + tamu.PhoneNumber +
				"\n" + tamu.Kategori,
		}
		go jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifsatpam, config.APIWAText)
		return "Selamat Kakak berhasil mengisi buku tamu dan terdaftar sebagai mitra pedagang athena"
	}
	//pemberitahuan ke satpam shift jaga
	notifsatpam := model.TextMessage{
		To:      satpam.Phonenumber,
		IsGroup: false,
		Messages: "*Kunjungan Tamu/Mitra lingkungan komplek*\nNama: " + tamu.Nama +
			"\nKontak: " + tamu.PhoneNumber +
			"\nKategori: " + tamu.Kategori +
			"\nTujuan: " + tamu.Tujuan +
			"\nBlokNomor: " + tamu.BlokRumah,
	}
	go jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifsatpam, config.APIWAText)
	//pemberitahuan ke penghuni
	notifpenghuni := model.TextMessage{
		To:      penghuniyangdikunjungi.Phonenumber,
		IsGroup: false,
		Messages: "*Kunjungan Tamu/Mitra lingkungan komplek*\nNama: " + tamu.Nama +
			"\nKontak: " + tamu.PhoneNumber +
			"\nKategori: " + tamu.Kategori +
			"\nTujuan: " + tamu.Tujuan +
			"\nBlokNomor: " + tamu.BlokRumah,
	}
	go jsonapi.PostStructWithToken[model.Response]("token", Profile.Token, notifpenghuni, config.APIWAText)

	return "Selamat Kakak berhasil mengisi buku tamu untuk kunjungan ke " + tamu.BlokRumah
}
