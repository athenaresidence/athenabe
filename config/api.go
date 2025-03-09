package config

import "os"

var PaperkaSecret = os.Getenv("PAPERKASECRET")

var APIWAText = "https://api.wa.my.id/api/v2/send/message/text"
