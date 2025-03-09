package config

import "os"

var PaperkaSecret = os.Getenv("PAPERKASECRET")

var APIWAText = "https://api.wa.my.id/api/v2/send/message/text"
var APIRegion = "https://asia-southeast2-awangga.cloudfunctions.net/petabackend/data/gis/lokasi"
