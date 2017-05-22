package networking

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"models"
	"net/http"
	"time"
)

// create type for consistent function signatures in routing map
type HandlerMethodType func(http.ResponseWriter, *http.Request, httprouter.Params)
type ControlMethodType func(string, httprouter.Handle)

func generateAPIEndPoint(fn HandlerMethodType) httprouter.Handle {
	return func(respWrtr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		reqStartTime := time.Now()
		defer func() {
			if err := recover(); err != nil {
				reqEndTime := time.Now()
				log.Printf("[%s] %q %v\n: panic: %+v\n", req.Method, req.URL.String(), reqEndTime.Sub(reqStartTime), err)
				models.WriteError(respWrtr, models.ErrInternalServer)
			}
		}()
		respWrtr.Header().Set("Sec-Websocket-Version", "13")
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		//respWrtr.Header().Set("Content-Type", "application/json")
		respWrtr.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, OPTION, HEAD, PATCH, PUT, POST, DELETE",
		)
		// A post should contain a request body
		if req.Method == "POST" && req.Body == nil {
			reqEndTime := time.Now()
			models.WriteError(respWrtr, models.ErrMissingPayload)
			log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), reqEndTime.Sub(reqStartTime))
			return
		}
		if req.URL.String() == "/api/docs" {
			reqEndTime := time.Now()
			json.NewEncoder(respWrtr).Encode(APIRouteMap)
			log.Printf("asdkjlfkjfklfklfkfjkffk")
			log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), reqEndTime.Sub(reqStartTime))
			return
		}
		fn(respWrtr, req, ps)
		reqEndTime := time.Now()
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), reqEndTime.Sub(reqStartTime))
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
		handler_method := Handles[end_point].(func(
			http.ResponseWriter,
			*http.Request,
			httprouter.Params,
		))

		log.Printf("GENERATING END POINT: ", ctrl_method, ": ", full_end_point)

		control_methods[ctrl_method](
			full_end_point,
			generateAPIEndPoint(handler_method),
		)
	}
	return router
}
