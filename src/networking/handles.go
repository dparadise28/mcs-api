package networking

import (
	"fmt"
	//"log"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Info(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI: %q\n", r.RequestURI)
	fmt.Fprintf(w, "URL: %#v\n", r.URL)
	fmt.Fprintf(w, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
	fmt.Fprintf(w, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	fmt.Fprintf(w, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(w, "\nHeaders:\n")
	r.Header.Write(w)
}

func Docs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routeDoc := make([]map[string]map[string]interface{}, len(APIRouteList))
	for index, route := range APIRouteList {
		for routeEndPoint, routeSpecs := range route {
			routeSpecs["api_method"] = ""
			routeDoc[index] = map[string]map[string]interface{}{routeEndPoint: routeSpecs}
		}
	}
	json.NewEncoder(w).Encode(routeDoc)
}
