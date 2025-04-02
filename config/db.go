package config

import (
	"os"

	"github.com/gocroot/mgdb"
)

// dhani
var mconn = mgdb.DBInfo{
	DBString: os.Getenv("MONGOSTRING"),
	DBName:   "athena",
}

var Mongoconn, ErrorMongoconn = mgdb.MongoConnect(mconn)

// paperka
