package kimseok

import "go.mongodb.org/mongo-driver/bson/primitive"

type Datasets struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Question string             `json:"question" bson:"question"`
	Answer   string             `json:"answer" bson:"answer"`
	Origin   string             `json:"origin" bson:"origin"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" `
	Username     string             `json:"username" bson:"username"`
	Email        string             `bson:"email,omitempty" json:"email,omitempty"`
	Password     string             `json:"password" bson:"password"`
	PasswordHash string             `json:"passwordhash" bson:"passwordhash"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	Private      string             `json:"private,omitempty" bson:"private,omitempty"`
	Public       string             `json:"public,omitempty" bson:"public,omitempty"`
}

type Secrets struct {
	SecretToken string `json:"secret_token" bson:"secret_token"`
}
