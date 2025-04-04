package helpdesk

import (
	"strings"

	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetHelpdeskName(WAPhoneNumber string, im model.WAMessage, MongoConn *mongo.Database, ModuleCollection string) (helpdesk Helpdesk, err error) {
	helpdesk, err = mgdb.GetOneDoc[Helpdesk](MongoConn, ModuleCollection, bson.M{"occupied": false})
	if err != nil {
		return
	}
	return
}

func IsMatch(str string, subs ...string) (bool, int) {

	matches := 0
	isCompleteMatch := true
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch, matches
}
