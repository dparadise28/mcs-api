package networking

import (
	"io/ioutil"
	"log"
	"net/http"
)

/*
func RedirectTarget(w http.ResponseWriter, req *http.Request, target string) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + *listen + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
*/
func LogExtIp(port string) {
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
