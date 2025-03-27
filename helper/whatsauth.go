package helper

import (
	"regexp"
	"strings"

	"github.com/gocroot/jsonapi"
	"github.com/gocroot/lite/helper/kimseok"
	"github.com/gocroot/lite/mod"
	"github.com/gocroot/lite/mod/presensi"
	"github.com/gocroot/lite/model"
	"github.com/gocroot/lite/module"
	"github.com/gocroot/mgdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WebHook(WAKeyword, WAPhoneNumber, WAAPIQRLogin, WAAPIMessage string, msg model.WAMessage, db *mongo.Database) (resp model.Response, err error) {
	if IsLoginRequest(msg, WAKeyword) { //untuk whatsauth request login
		resp, err = HandlerQRLogin(msg, WAKeyword, WAPhoneNumber, db, WAAPIQRLogin)
	} else { //untuk membalas pesan masuk
		resp, err = HandlerIncomingMessage(msg, WAPhoneNumber, db, WAAPIMessage)
	}
	return
}

func RefreshToken(dt *model.WebHook, WAPhoneNumber, WAAPIGetToken string, db *mongo.Database) (res *mongo.UpdateResult, err error) {
	profile, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	var resp model.User
	if profile.Token != "" {
		_, resp, err = jsonapi.PostStructWithToken[model.User]("Token", profile.Token, dt, WAAPIGetToken)
		if err != nil {
			return
		}
		profile.Phonenumber = resp.PhoneNumber
		profile.Token = resp.Token
		res, err = mgdb.ReplaceOneDoc(db, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
		if err != nil {
			return
		}
	}
	return
}

func IsLoginRequest(msg model.WAMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) // && msg.From_link
}

func GetUUID(msg model.WAMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg model.WAMessage, WAKeyword string, WAPhoneNumber string, db *mongo.Database, WAAPIQRLogin string) (resp model.Response, err error) {
	dt := &model.WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Aliasname:   msg.Alias_name,
		Delay:       msg.From_link_delay,
	}
	structtoken, err := GetAppProfile(WAPhoneNumber, db)
	if err != nil {
		return
	}
	_, resp, err = jsonapi.PostStructWithToken[model.Response]("Token", structtoken.Token, dt, WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg model.WAMessage, WAPhoneNumber string, db *mongo.Database, WAAPIMessage string) (resp model.Response, err error) {
	_, bukanbot := GetAppProfile(msg.Phone_number, db) //cek apakah nomor adalah bot
	if bukanbot != nil {                               //jika tidak terdapat di profile
		var profile model.Profile
		profile, err = GetAppProfile(WAPhoneNumber, db)
		if err != nil {
			return
		}
		msg.Message = NormalizeHiddenChar(msg.Message)
		module.NormalizeAndTypoCorrection(&msg.Message, db, "typo")
		modname, group, personal := module.GetModuleName(WAPhoneNumber, msg, db, "module")
		var msgstr string
		var isgrup bool
		if msg.Chat_server != "g.us" { //chat personal
			if personal && modname != "" {
				msgstr = mod.Caller(profile, modname, msg, db)
			} else {
				msgstr = kimseok.GetMessage(profile, msg, profile.Botname, db)
			}

		} else if strings.Contains(strings.ToLower(msg.Message), profile.Triggerword) { //chat group
			isgrup = true
			if group && modname != "" {
				msgstr = mod.Caller(profile, modname, msg, db)
			} else {
				msgstr = kimseok.GetMessage(profile, msg, profile.Botname, db)
			}
		}
		dt := &model.TextMessage{
			To:       msg.Chat_number,
			IsGroup:  isgrup,
			Messages: msgstr,
		}
		_, resp, err = jsonapi.PostStructWithToken[model.Response]("Token", profile.Token, dt, WAAPIMessage)
		if err != nil {
			return
		}

	}
	return
}

func GetRandomReplyFromMongo(msg model.WAMessage, botname string, db *mongo.Database) string {
	rply, err := mgdb.GetRandomDoc[model.Reply](db, "reply", 1)
	if err != nil {
		return "Koneksi Database Gagal: " + err.Error()
	}
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", botname)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func GetMessageFromKimseokgis(msg model.WAMessage, botname string, db *mongo.Database) string {
	conf, err := mgdb.GetOneDoc[presensi.Config](db, "config", bson.M{"phonenumber": "62895601060000"})
	if err != nil {
		return "data config gagal di ambil " + err.Error()
	}
	dt := model.Requests{Messages: msg.Message}
	_, msgreply, err := jsonapi.PostStructWithToken[model.Chats]("secret", conf.LeaflySecret, dt, conf.KimseokgisURL)
	if err != nil {
		return "Wah kak kayaknya salah pas ngepost reply deh " + err.Error()
	}
	return msgreply.Responses
}

func GetAppProfile(phonenumber string, db *mongo.Database) (apitoken model.Profile, err error) {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, err = mgdb.GetOneDoc[model.Profile](db, "profile", filter)

	return
}

func NormalizeHiddenChar(text string) string {
	return removeZeroWidthSpaces(removeInvisibleChars(text))
}

func removeZeroWidthSpaces(text string) string {
	// Create a regular expression to match specific zero-width characters
	re := regexp.MustCompile(`\p{Cf}`)

	// Replace all matches with an empty string
	return re.ReplaceAllString(text, "")
}

func removeInvisibleChars(text string) string {
	// Create a regular expression to match invisible characters
	re := regexp.MustCompile(`\p{C}`)

	// Replace all matches with an empty string
	return re.ReplaceAllString(text, "")
}
