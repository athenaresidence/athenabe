package bukutamu

import (
	"regexp"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
)

func BukuTamu(Profile model.Profile, Pesan model.WAMessage) (reply string) {
	form := "https://wa.me/6281776900300?text=%2A8ukuT4muAth3na%2A+%0A%3E+Asal+%3A+%0A%3E+Tujuan+%3A+%0A%3E+NoPol+%3A+"
	tamu := ParsePesanFlexible(Pesan.Message)
	tamu.PhoneNumber = Pesan.Phone_number
	tamu.Nama = Pesan.Alias_name
	//cek dahulu isian asal dan nopol
	if tamu.Kendaraan == "" || tamu.Kategori == "" {
		return "Mohon maaf asal dan nopol harus diisi kak, silahkan coba kembali:\n" + form
	}
	//mengambil data satpam yang masuk shift aktif terakhir
	satpam, err := mgdb.GetOneLatestDoc[model.Satpam](config.Mongoconn, "selfie", bson.M{})
	if err != nil {
		return "Mohon maaf saat ini belum ada satpam yang menjaga pos"
	}

	//pemberitahuan ke grup untuk ada yang jualan masuk ke komplek
	re := regexp.MustCompile(`\s+`)
	cleanTujuan := re.ReplaceAllString(tamu.Tujuan, "")
	if cleanTujuan == "niaga" {
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
		mgdb.InsertOneDoc(config.Mongoconn, "bukutamu", tamu)
		return "Selamat Kakak berhasil mengisi buku tamu dan terdaftar sebagai mitra pedagang athena"
	}
	//mengambil data penghuni yang dituju
	penghuniyangdikunjungi, err := mgdb.GetOneDoc[model.Penghuni](config.Mongoconn, "penghuni", bson.M{"nomorrumah": tamu.BlokRumah})
	if err != nil && tamu.Tujuan != "niaga" {
		return "Mohon maaf tujuan yang kakak input tidak valid:\n_" + tamu.Tujuan + "_" + "\n silahkan coba kembali:\n" + form
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
	mgdb.InsertOneDoc(config.Mongoconn, "bukutamu", tamu)
	return "Selamat Kakak berhasil mengisi buku tamu untuk kunjungan ke " + tamu.BlokRumah
}
