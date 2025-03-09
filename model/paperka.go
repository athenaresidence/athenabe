package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResellerPaperka struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Alamat      string             `json:"alamat" bson:"alamat"`
	Kelurahan   string             `json:"kelurahan" bson:"kelurahan"`
	Kecamatan   string             `json:"kecamatan" bson:"kecamatan"`
	Kota        string             `json:"kota" bson:"kota"`
	Provinsi    string             `json:"provinsi" bson:"provinsi"`
	Email       string             `json:"email" bson:"email"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
}

type Session struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"userid"`
	CreatedAt time.Time `bson:"createdAt"`
}
