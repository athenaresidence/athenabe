package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wag struct {
	Phonenumber string `json:"phonenumber" bson:"phonenumber"`
}

type Penghuni struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Nomorrumah  string             `json:"nomorrumah,omitempty" bson:"nomorrumah,omitempty"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
	Kendaraan   []Kendaraan        `json:"kendaraan,omitempty" bson:"kendaraan,omitempty"`
}

type Pengantar struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Nomorrumah  string             `json:"nomorrumah,omitempty" bson:"nomorrumah,omitempty"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
	Kendaraan   []Kendaraan        `json:"kendaraan,omitempty" bson:"kendaraan,omitempty"`
}

type Tamu struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Nomorrumah  string             `json:"nomorrumah,omitempty" bson:"nomorrumah,omitempty"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
	Kendaraan   []Kendaraan        `json:"kendaraan,omitempty" bson:"kendaraan,omitempty"`
}

type Kendaraan struct {
	Merk      string `json:"merk,omitempty" bson:"merk,omitempty"`
	PlatNomor string `json:"platnomor,omitempty" bson:"platnomor,omitempty"`
	Warna     string `json:"warna,omitempty" bson:"warna,omitempty"`
}

type Satpam struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Nama        string             `json:"nama" bson:"nama"`
	Phonenumber string             `json:"phonenumber" bson:"phonenumber"`
}

type JadwalPos struct {
	Shift string
}

type Header struct {
	Secret string `reqHeader:"secret,omitempty"` //whatsauth ke webhook
	Token  string `reqHeader:"token,omitempty"`  //webhook ke whatsauth kirim pesan
}

type Profile struct {
	Token            string `bson:"token"`
	Phonenumber      string `bson:"phonenumber"`
	AdminPhonenumber string `bson:"adminphonenumber"`
	Secret           string `bson:"secret"`
	URL              string `bson:"url"`
	URLAPIText       string `bson:"urlapitext"`
	URLAPIImage      string `bson:"urlapiimage"`
	URLAPIDoc        string `bson:"urlapidoc"`
	URLQRLogin       string `bson:"urlqrlogin"`
	QRKeyword        string `bson:"qrkeyword"`
	PublicKey        string `bson:"publickey"`
	Botname          string `bson:"botname"`
	Triggerword      string `bson:"triggerword"`
	TelegramToken    string `bson:"telegramtoken"`
	TelegramName     string `bson:"telegramname"`
}
