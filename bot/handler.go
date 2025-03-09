package bot

import (
	"fmt"
	"strings"

	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/helper/atapi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlerPesan(msg model.WAMessage, profile model.Profile) (reply string) {
	user, err := mgdb.GetOneDoc[model.UserResellerPaperka](config.Mongoconnpaperka, "user", bson.M{"phonenumber": msg.Phone_number})
	var userbelumterdaftar bool
	if err == mongo.ErrNoDocuments {
		userbelumterdaftar = true
	}
	if msg.Latitude != 0.0 {
		loc := model.LongLat{
			Longitude: msg.Longitude,
			Latitude:  msg.Latitude,
		}
		_, res, _ := atapi.PostStructWithToken[model.Region]("Login", profile.Token, loc, config.APIRegion)
		reply = fmt.Sprintf("Lokasi kak %s di:\n%s %s %s %s\nLat:%.2f Long:%.2f", msg.Alias_name, res.Village, res.SubDistrict, res.District, res.Province, msg.Latitude, msg.Longitude)
		user.Kelurahan = res.Village
		user.Kecamatan = res.SubDistrict
		user.Kota = res.District
		user.Provinsi = res.Province
		if userbelumterdaftar {
			user.Nama = msg.Alias_name
			user.Phonenumber = msg.Phone_number
			idinsert, err := mgdb.InsertOneDoc(config.Mongoconnpaperka, "user", user)
			if err != nil {
				reply += err.Error()
			} else {
				reply += idinsert.String()
			}
		} else {
			updateFields := bson.M{
				"kelurahan": res.Village,
				"kecamatan": res.SubDistrict,
				"kota":      res.District,
				"provinsi":  res.Province,
			}
			mgdb.UpdateOneDoc(config.Mongoconnpaperka, "user", bson.M{"phonenumber": msg.Phone_number}, updateFields)
		}
		return
	}
	if strings.Contains(msg.Message, "alamatpengirimanpaperka:") {
		result := strings.SplitN(msg.Message, ":", 2)
		alamat := strings.TrimSpace(result[1])
		reply = "alamat pengiriman kakak :\n" + alamat
		user.Alamat = alamat
		if userbelumterdaftar {
			user.Nama = msg.Alias_name
			user.Phonenumber = msg.Phone_number
			mgdb.InsertOneDoc(config.Mongoconnpaperka, "user", user)
		} else {
			updateFields := bson.M{
				"alamat": alamat,
			}
			mgdb.UpdateOneDoc(config.Mongoconnpaperka, "user", bson.M{"phonenumber": msg.Phone_number}, updateFields)
		}
		return
	}
	if userbelumterdaftar {
		reply = "Selamat datang kak " + msg.Alias_name
		reply += "\nSilahkan share lock lokasi pengiriman dulu kak"
		user := model.UserResellerPaperka{
			Nama:        msg.Alias_name,
			Phonenumber: msg.Phone_number,
		}
		mgdb.InsertOneDoc(config.Mongoconnpaperka, "user", user)
		return
	} else {
		if user.Provinsi == "" {
			reply = "Selamat datang kak " + msg.Alias_name
			reply += "\nKakak belum share loc alamat pengiriman silahkan share lock lokasi pengiriman kak"
		} else if user.Alamat == "" {
			reply = "Selamat datang kak " + msg.Alias_name
			reply += "\nKakak belum mengisi alamat silahkan mengisi alamat pengiriman dengan mengetikkan *alamatpengirimanpaperka:* di depan alamat atau klik saja disini https://wa.me/628112109691?text=alamatpengirimanpaperka:%0A"
		} else {
			reply = UserTerdaftar(user)
		}
	}
	return

}

func UserTerdaftar(user model.UserResellerPaperka) (reply string) {
	_, err := mgdb.GetOneDoc[model.Session](config.Mongoconnpaperka, "session", bson.M{"userid": user.Phonenumber})
	if err == mongo.ErrNoDocuments {
		reply = "Selamat datang kak " + user.Nama
		reply += "\nAlamat pengiriman:\n" + user.Alamat + "\n" + user.Kelurahan + "," + user.Kecamatan + "," + user.Kota + "," + user.Provinsi
		reply = "\n\nMohon tunggu sebentar, mimin sebentar lagi akan membalas"
		return
	}
	return
}
