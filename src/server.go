package main

import (
	"crypto/tls"
	"db"
	"flag"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"models"
	"net/http"
	"networking"
	"os"
	"path/filepath"
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

var m autocert.Manager

func RedirectHTTPS(w http.ResponseWriter, req *http.Request) {
	// TODO: move out of base server and into networking utils

	// remove/add non default ports from req.Host
	target := "https://" + models.DOMAIN + req.URL.Path
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
	flag.StringVar(&tools.SLACK_TOKEN, "slack_token", "", "Your slack token for notifications")
	flag.StringVar(&models.UI_DIR_PATH, "ui_dir_path", "", "Your slack token for notifications")
	flag.StringVar(&models.DOMAIN, "domain", "", "Your servers domain")
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
		len(models.DOMAIN) == 0,
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
}

func setupTLS() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't find working directory to locate or save certificates.")
	}

	cache := autocert.DirCache(filepath.Join(pwd, "system", "tls", "certs"))
	if _, err := os.Stat(string(cache)); os.IsNotExist(err) {
		err := os.MkdirAll(string(cache), os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatalln("Couldn't create cert directory at", cache)
		}
	}

	// get host/domain and email from Config to use for TLS request to Let's encryption.
	// we will fail fatally if either are not found since Let's Encrypt will rate-limit
	// and sending incomplete requests is wasteful and guaranteed to fail its check
	host := models.DOMAIN
	m = autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       cache,
		HostPolicy:  autocert.HostWhitelist(string(host)),
		RenewBefore: time.Hour * 24 * 30,
		Email:       string(tools.PlatformEmail),
	}

}

func main() {
	/*f, err := os.OpenFile("mcs-api.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()*/
	if port == "443" {
		rw := tools.NewRotateWriter("mcs-api.log")
		log.SetOutput(rw)
		defer rw.Close()
	}

	db.InitSession()
	db.InitIndicies()
	setupTLS()
	defer db.Database.Session.Close()
	listen := ":" + port
	models.JWT_SIGNATURE = models.StripeSK

	// redirect every http request to https
	go http.ListenAndServe(":80", http.HandlerFunc(RedirectHTTPS))
	// heartbeat for aws and general monitoring
	go http.ListenAndServe(":8080", http.HandlerFunc(networking.InfoHandler))

	log.Println("\n\n-----Starting Endpoints\n")
	server := &http.Server{
		Addr:           listen,
		Handler:        http.Handler(networking.ServeEndPoints()),
		ReadTimeout:    6 * time.Second,
		WriteTimeout:   6 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      &tls.Config{GetCertificate: m.GetCertificate},
	}
	networking.LogExtIp(server.Addr)
	if listen != ":443" {
		log.Fatal(server.ListenAndServe())
		//log.Fatal(server.ListenAndServeTLS("src/certs/server/cert.pem", "src/certs/server/key.pem"))
	} else {
		log.Println("Starting TLS")
		log.Fatal(server.ListenAndServeTLS("", ""))
	}
}
