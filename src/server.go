package main

import (
	//"github.com/julienschmidt/httprouter"
	"crypto/tls"
	"flag"
	//"io/ioutil"
	"db"
	"log"
	"net/http"
	"networking"
	//"os"
	"time"
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

func InitServer() {
	flag.StringVar(&port, "port", "", "Port input by user when starting server")
	flag.StringVar(&db.AuthDatabase, "db_name", "", "Name of the db you would like to connect to")
	flag.StringVar(&db.MongoDBUri, "mongo_db_uri", "",
		"example: mongodb://<dbuser>:<dbpassword>@<dbhost1>,<dbhost2>,...:<port>/<dbname>")
	flag.Parse()
	if len(port) == 0 || len(db.MongoDBUri) == 0 || len(db.AuthDatabase) == 0 {
		panic(`
			Must provide a port to start the server on, a db uri to connect to
			and a db name to use.

			ex: go run server.go --port=443 --db_name=test-db --mongo_db_uri=mongodb://...
		`)
	}
}

func main() {
	InitServer()
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
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}
	networking.LogExtIp(server.Addr)
	if listen != ":443" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("Starting TLS")
		log.Fatal(server.ListenAndServeTLS("certs/server/cert.pem", "certs/server/key.pem"))
	}
}
