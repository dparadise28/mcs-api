package config

import (
	"time"
)

var (
	// server
	SSL        = false
	Port       = ""
	SSLCertTTL = time.Hour * 24 * 3

	// mongo
	MongoDBUri            = ""
	AuthDatabase          = ""
	StripeSK              = ""
	EmailPassword         = ""
	PlatformEmail         = ""
	AWS_SECRET_ACCESS_KEY = ""
	AWS_ACCESS_KEY_ID     = ""
	AWS_REGION            = ""
	AWS_S3_BUCKET_NAME    = ""
	AWS_S3_BUCKET_KEY     = ""
	DOMAIN                = ""
)
