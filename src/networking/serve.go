package networking

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// create type for consistent function signatures in routing map
type HandlerMethodType func(http.ResponseWriter, *http.Request, httprouter.Params)
type ControlMethodType func(string, httprouter.Handle)

func generateAPIEndPoint(fn HandlerMethodType) httprouter.Handle {
	return func(respWrtr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		respWrtr.Header().Set("Sec-Websocket-Version", "13")
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		respWrtr.Header().Set("Content-Type", "application/json")
		respWrtr.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, OPTION, HEAD, PATCH, PUT, POST, DELETE",
		)
		log.Printf(
			"Request for %s (Accept-Encoding: %s)",
			req.URL.Path,
			req.Header.Get("Accept-Encoding"),
		)
		fn(respWrtr, req, ps)
	}
}

func ServeEndPoints() *httprouter.Router {
	router := httprouter.New()
	control_methods := map[string]ControlMethodType{
		// Not much else is needed for now
		"GET":  router.GET,
		"POST": router.POST,
	}
	for end_point, api_end_point := range APIRouteMap {
		ctrl_method := api_end_point["control_method"].(string)
		full_end_point := "/api" + end_point
		handler_method := api_end_point["handler_method"].(func(
			http.ResponseWriter,
			*http.Request,
			httprouter.Params,
		))

		log.Println("GENERATING END POINT: ", ctrl_method, ": ", full_end_point)

		control_methods[ctrl_method](
			full_end_point,
			generateAPIEndPoint(handler_method),
		)
	}
	return router
}
