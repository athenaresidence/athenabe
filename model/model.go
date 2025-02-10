package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wag struct {
	Phonenumber string `json:"phonenumber" bson:"phonenumber"`
}

type Penghuni struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Nomorrumah  string             `json:"nomorrumah" bson:"nomorrumah"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
	Kendaraan   []Kendaraan        `json:"kendaraan" bson:"kendaraan"`
}

type Kendaraan struct {
	Merk      string
	PlatNomor string
	Warna     string
}

type Satpam struct {
	Nama        string
	Phonenumber string
}

type JadwalPos struct {
	Shift string
}
