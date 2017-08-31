package main

import (
	"crypto/tls"
	"db"
	"flag"
	"log"
	"models"
	"net/http"
	"networking"
	"time"
	"tools"
)

var (
	ssl            = false
	port           = ""
	certPath       = ""
	sslCrtFileName = ""
	sslKeyFileName = ""
)

func RedirectHTTPS(w http.ResponseWriter, req *http.Request) {
	// TODO: move out of base server and into networking utils

	// remove/add non default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)

}

func init() {
	//general
	flag.StringVar(&port, "port", "", "Port input by user when starting server")
	// mongo
	flag.StringVar(&db.MongoDBUri, "mongo_db_uri", "", "example: mongodb://<dbuser>:<dbpassword>@<dbhost1>,<dbhost2>,...:<port>/<dbname>")
	flag.StringVar(&db.AuthDatabase, "db_name", "", "Name of the db you would like to connect to")
	// stripe payments
	flag.StringVar(&models.StripeSK, "stripe_sk", "", "Your stripe secret key")
	// emails
	flag.StringVar(&tools.EmailPassword, "platform_email_pw", "", "The password to your platforms email account")
	flag.StringVar(&tools.PlatformEmail, "platform_email_addr", "", "Your platforms email address (this is used for things like sending confirmation emails)")
	// aws
	flag.StringVar(&models.AWS_SECRET_ACCESS_KEY, "aws_secret_access_key", "", "Your aws secret access key")
	flag.StringVar(&models.AWS_ACCESS_KEY_ID, "aws_access_key_id", "", "Your aws access key id")
	flag.StringVar(&models.AWS_REGION, "aws_region", "", "Your aws region")
	flag.StringVar(&models.AWS_S3_BUCKET_NAME, "aws_s3_bucket_name", "", "Your aws s3 bucket name")
	flag.StringVar(&models.AWS_S3_BUCKET_KEY, "aws_s3_bucket_key", "", "Your aws s3 bucket key")
	flag.Parse()
	for _, condition := range []bool{
		len(port) == 0,
		len(db.MongoDBUri) == 0,
		len(db.AuthDatabase) == 0,
		len(models.StripeSK) == 0,
		len(tools.EmailPassword) == 0,
		len(tools.PlatformEmail) == 0,
		len(models.AWS_SECRET_ACCESS_KEY) == 0,
		len(models.AWS_ACCESS_KEY_ID) == 0,
		len(models.AWS_REGION) == 0,
		len(models.AWS_S3_BUCKET_NAME) == 0,
		len(models.AWS_S3_BUCKET_KEY) == 0,
	} {
		if condition {
			panic(`
			Must provide a port to start the server on, a db uri to connect to
			and a db name to use.

			ex: 
				go run server.go --port=443 \
								 --db_name=test-db \
								 --mongo_db_uri=mongodb://... \
								 --platform_email_pw=pa530rd \
								 --platform_email_addr=example@domain.com \
								 --stripe_sk=s3cR3tKey
		`)
		}
	}
	log.Println()
}

func main() {
	db.InitSession()
	db.InitIndicies()
	defer db.Database.Session.Close()
	listen := ":" + port //os.Args[1])

	// redirect every http request to https
	go http.ListenAndServe(":80", http.HandlerFunc(RedirectHTTPS))

	log.Println("\n\n-----Starting Endpoints\n")
	server := &http.Server{
		Addr:         listen,
		Handler:      http.Handler(networking.ServeEndPoints()),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		// MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			NextProtos:               []string{"h2", "h2-14"},
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		},
	}
	networking.LogExtIp(server.Addr)
	if listen != ":443" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("Starting TLS")
		// log.Fatal(server.ListenAndServeTLS("certs/server/cert.pem", "certs/server/key.pem"))
		log.Fatal(server.ListenAndServeTLS("../mcs_ssl/mycorner_store/ssl-bundle.crt", "../mcs_ssl/mycorner_store/mycorner_store.pkey"))
	}
}
