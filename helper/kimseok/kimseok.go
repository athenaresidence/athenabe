package kimseok

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gocroot/lite/mod/helpdesk"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/mgdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCursorFromRegex(db *mongo.Database, patttern string) (cursor *mongo.Cursor, err error) {
	filter := bson.M{"question": primitive.Regex{Pattern: patttern, Options: "i"}}
	cursor, err = db.Collection("qna").Find(context.TODO(), filter)
	return
}

func GetCursorFromString(db *mongo.Database, question string) (cursor *mongo.Cursor, err error) {
	filter := bson.M{"question": question}
	cursor, err = db.Collection("qna").Find(context.TODO(), filter)
	return
}

func GetRandomFromQnASlice(qnas []Datasets) Datasets {
	// Inisialisasi sumber acak dengan waktu saat ini
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Pilih elemen acak dari slice
	randomIndex := rng.Intn(len(qnas))
	return qnas[randomIndex]
}

func GetQnAfromSliceWithJaro(q string, qnas []Datasets) (dt Datasets) {
	var score float64
	for _, qna := range qnas {
		//fmt.Println(data)
		str2 := qna.Question
		scorex := jaroWinkler(q, str2)
		if score < scorex {
			dt = qna
			score = scorex
		}

	}
	return

}

// balasan jika tidak ditemukan key word
func GetMessage(Profile model.Profile, msg model.WAMessage, botname string, db *mongo.Database) string {
	//check apakah ada permintaan operator masuk
	reply, err := helpdesk.PenugasanOperator(Profile, msg, db)
	if err != nil {
		return err.Error()
	}
	//jika tidak ada di db komplain lanjut ke selanjutnya
	if reply == "" {
		dt, err := QueriesDataRegexpALL(db, msg.Message)
		if err != nil {
			return err.Error()
		}
		reply = strings.TrimSpace(dt.Answer)

	}
	return reply
}

func QueriesDataRegexpALL(db *mongo.Database, queries string) (dest Datasets, err error) {
	//kata akhiran imbuhan
	//queries = SeparateSuffixMu(queries)

	//ubah ke kata dasar
	queries = Stemmer(queries)
	filter := bson.M{"question": queries}
	qnas, err := mgdb.GetAllDoc[[]Datasets](db, "qna", filter)
	if err != nil {
		return
	}
	if len(qnas) > 0 {
		dest = GetRandomFromQnASlice(qnas)
		return
	}

	wordsdepan := strings.Fields(queries) // Tokenization
	wordsbelakang := wordsdepan
	//pencarian dengan pengurangan kata dari belakang
	for len(wordsdepan) > 0 || len(wordsbelakang) > 0 {
		// Join remaining elements back into a string for wordsdepan
		filter := bson.M{"question": primitive.Regex{Pattern: strings.Join(wordsdepan, " "), Options: "i"}}
		qnas, err = mgdb.GetAllDoc[[]Datasets](db, "qna", filter)
		if err != nil {
			return
		} else if len(qnas) > 0 {
			dest = GetQnAfromSliceWithJaro(queries, qnas)
			return
		}
		// Join remaining elements back into a string for wordsbelakang
		filter = bson.M{"question": primitive.Regex{Pattern: strings.Join(wordsbelakang, " "), Options: "i"}}
		qnas, err = mgdb.GetAllDoc[[]Datasets](db, "qna", filter)
		if err != nil {
			return
		} else if len(qnas) > 0 {
			dest = GetQnAfromSliceWithJaro(queries, qnas)
			return
		}
		// Remove the last element
		wordsdepan = wordsdepan[:len(wordsdepan)-1]
		// remove element pertama
		wordsbelakang = wordsbelakang[1:]
	}

	return
}
