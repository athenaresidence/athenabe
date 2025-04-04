package strava

import (
	"errors"
	"net/http"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getConfigByPhone(db *mongo.Database, profilePhone string) (*Config, error) {
	conf, err := mgdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": profilePhone})
	if err != nil {
		return nil, errors.New("kesalahan dalam pengambilan config di database: " + err.Error())
	}
	return &conf, nil
}

func postToDomyikado(secret, url string, data map[string]interface{}) error {
	statuscode, httpresp, err := jsonapi.PostStructWithToken[model.Response]("secret", secret, data, url)
	if err != nil {
		return errors.New("akses ke endpoint domyikado gagal: " + err.Error())
	}

	if statuscode != http.StatusOK {
		return errors.New("salah posting endpoint domyikado: " + httpresp.Response + "\ninfo: " + httpresp.Info)
	}

	return nil
}
