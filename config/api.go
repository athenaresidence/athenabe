package config

import "os"

var PaperkaSecret = os.Getenv("PAPERKASECRET")
