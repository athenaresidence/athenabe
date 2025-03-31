package satpam

import "go.mongodb.org/mongo-driver/bson/primitive"

type LaporanBulanan struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"token"`
	Message string             `json:"message" bson:"message"`
}
