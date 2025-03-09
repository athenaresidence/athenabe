package bot

import (
	"fmt"
	"strings"

	"github.com/gocroot/lite/config"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandlerPesan(msg model.WAMessage, profile model.Profile) (reply string) {
	user, err := mgdb.GetOneDoc[model.UserResellerPaperka](config.Mongoconnpaperka, "user", bson.M{"phonenumber": msg.Phone_number})
	if msg.Latitude != 0.0 {
		reply = fmt.Sprintf("Lokasi kak %s di:\n Lat:%.2f Long:%.2f", msg.Alias_name, msg.Latitude, msg.Longitude)
		return
	}
	if strings.Contains(msg.Message, "alamatpengirimanpaperka:") {
		result := strings.SplitN(msg.Message, ":", 2)
		alamat := strings.TrimSpace(result[1])
		reply = "alamat pengiriman kakak :\n" + alamat
		user.Alamat = alamat
		mgdb.ReplaceOneDoc(config.Mongoconnpaperka, "user", bson.M{"phonenumber": msg.Phone_number}, user)
		return
	}
	if err == mongo.ErrNoDocuments {
		reply = "Selamat datang kak " + msg.Alias_name
		reply += "\nSilahkan share lock lokasi pengiriman dulu kak"
		user := model.UserResellerPaperka{
			Nama:        msg.Alias_name,
			Phonenumber: msg.Phone_number,
		}
		mgdb.InsertOneDoc(config.Mongoconnpaperka, "user", user)
		return
	}
	return

}

func PendaftaranUser() {

}
