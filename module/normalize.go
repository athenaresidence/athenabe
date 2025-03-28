package module

import (
	"regexp"
	"strings"

	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NormalizeAndTypoCorrection(message *string, MongoConn *mongo.Database, TypoCollection string) {
	typos, _ := mgdb.GetAllDoc[[]Typo](MongoConn, TypoCollection, bson.M{})
	for _, typo := range typos {
		re := regexp.MustCompile(`(?i)` + typo.From + ``)
		*message = re.ReplaceAllString(*message, typo.To)
	}
	//merubah ke huruf kecil semua
	*message = strings.ToLower(*message)
	*message = strings.TrimSpace(*message)
	*message = cleanString(*message)
}

func cleanString(input string) string {
	// Trim leading & trailing spaces/tabs
	cleaned := strings.TrimSpace(input)

	// Replace multiple spaces/tabs with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	return cleaned
}
