package main

import (
	//"github.com/julienschmidt/httprouter"
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"networking"
	"os"
	"time"
)

var (
	listen = flag.String("listen", ":", "Port to listen on")
	port   = flag.String("port", "", "Port input by user when starting server")
	//config = flag.String("config", "", "Config file")
)

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + *listen + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func log_ext_ip(port string) {
	// lets log our external ip for easy access
	resp, err := http.Get("http://myexternalip.com/raw")
	if err == nil {
		extip, extipErr := ioutil.ReadAll(resp.Body)
		if extipErr == nil {
			log.Println("Setting Server Address", string(extip[:len(extip)-1])+port)
		} else {
			log.Println("\n\nTrouble Parsing external ip\n\nSetting Server Address", port)
		}
	} else {
		// shouldnt stop the server from starting
		log.Println(err.Error())
		log.Println("\n\nTrouble Retreiving external ip\n\nSetting Server Address", port)
	}
	resp.Body.Close()
}

func main() {
	flag.Parse()
	*listen += string(os.Args[1])

	// redirect every http request to https
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	log.Println("\n\n-----Starting Endpoints\n")
	server := &http.Server{
		Addr:         *listen,
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
	log_ext_ip(server.Addr)

	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual

	if *listen != ":443" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("Starting TLS")
		log.Fatal(server.ListenAndServeTLS("certs/server/cert.pem", "certs/server/key.pem"))
	}
}
