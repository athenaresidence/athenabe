package config

import "os"

var PaperkaSecret = os.Getenv("PAPERKASECRET")
var AthenaBotNumber = os.Getenv("BOTNUMBER")

var APIWAText = "https://api.wa.my.id/api/v2/send/message/text"
var APIWAIMG = "https://api.wa.my.id/api/send/message/image"
var APIRegion = "https://asia-southeast2-awangga.cloudfunctions.net/petabackend/data/gis/lokasi"
var APISignUp = "https://api.wa.my.id/api/signup"
var WAAPIQRLogin string = "https://api.wa.my.id/api/whatsauth/request"
