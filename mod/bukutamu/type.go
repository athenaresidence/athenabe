package bukutamu

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tamu struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nama        string             `bson:"nama,omitempty"`
	PhoneNumber string             `bson:"phonenumber,omitempty"`
	Kategori    string             `bson:"kategori,omitempty"`  //tamu atau nama mitra
	Tujuan      string             `bson:"tujuan,omitempty"`    //kosongkan untuk penjual
	BlokRumah   string             `bson:"blokrumah,omitempty"` // <- hasil parsing, misalnya "C12"
	Kendaraan   string             `bson:"kendaraan,omitempty"`
}
