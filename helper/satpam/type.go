package satpam

import "go.mongodb.org/mongo-driver/bson/primitive"

type LaporanBulanan struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Message string             `json:"message" bson:"message"`
}

type RekapRating struct {
	PhoneNumber     string      `json:"phonenumber"`
	TotalRating     int         `json:"total_rating"`
	AverageRating   float64     `json:"average_rating"`
	JumlahPerRating map[int]int `json:"jumlah_per_rating"`
}
