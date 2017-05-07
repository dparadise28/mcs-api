package api

import (
	"file_server"
	"log"
	"net/http"
)

func main() {
	var srv http.Server
	srv.Addr = ":8000"

	http.HandleFunc("/", file_server.ManageFile)
	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
