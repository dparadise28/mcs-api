package networking

import (
	//"github.com/didip/tollbooth"
	//"github.com/didip/tollbooth/thirdparty/tollbooth_httprouter"
	//"reflect"
	"github.com/julienschmidt/httprouter"
	"log"
	"models"
	"net/http"
	"time"
)

// create type for consistent function signatures in routing map
type HandlerMethodType func(http.ResponseWriter, *http.Request, httprouter.Params)
type ControlMethodType func(string, httprouter.Handle)

func generateAPIEndPoint(fn HandlerMethodType, fullEndPoint string) httprouter.Handle {
	return func(respWrtr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		reqStartTime := time.Now()
		defer func() {
			if err := recover(); err != nil {
				reqEndTime := time.Now()
				log.Printf("[%s] %s %s %v\n: panic: %+v\n",
					req.Method,
					req.URL.String(),
					fullEndPoint,
					reqEndTime.Sub(reqStartTime),
					err,
				)
				models.WriteError(respWrtr, models.ErrInternalServer)
			}
		}()

		// found in network/security
		SetBaseHeaders(respWrtr, req)
		if req.Method == "OPTIONS" {
			return
		}

		// A post should contain a request body
		if req.Method == "POST" && req.Body == nil {
			reqEndTime := time.Now()
			models.WriteError(respWrtr, models.ErrMissingPayload)
			log.Printf("[%s] %s %s %v\n", req.Method, req.URL.String(), fullEndPoint, reqEndTime.Sub(reqStartTime))
			return
		}
		ep := fullEndPoint[4:]
		if len(APIRouteMap[ep]["authenticate"].([]string)) != 0 {
			if valid, errJSON := ValidatedToken(respWrtr, req, ps, ep); !valid {
				models.WriteError(respWrtr, errJSON)
				reqEndTime := time.Now()
				log.Printf("[%s] :UNAUTHORIZED: %s %q %v\n", req.Method, req.URL.String(), fullEndPoint, reqEndTime.Sub(reqStartTime))
				return
			}
		}
		fn(respWrtr, req, ps)
		reqEndTime := time.Now()
		log.Printf("[%s] %s %q %v\n", req.Method, req.URL.String(), fullEndPoint, reqEndTime.Sub(reqStartTime))
	}
}

func ServeEndPoints() *httprouter.Router {
	router := httprouter.New()
	control_methods := map[string]ControlMethodType{
		// Not much else is needed for now (preflight options are automatic for all calls)
		"GET":  router.GET,
		"POST": router.POST,
	}
	for end_point, api_end_point := range APIRouteMap {
		ctrl_method := api_end_point["control_method"].(string)
		fullEndPoint := "/api" + end_point
		handler_method := APIRouteMap[end_point]["api_method"].(func(
			http.ResponseWriter,
			*http.Request,
			httprouter.Params,
		))

		log.Printf("GENERATING END POINT: ", ctrl_method, ": ", fullEndPoint)
		/*var rps int64
		if api_end_point["max_rps"] != nil {
			rps = int64(api_end_point["max_rps"].(int))
		} else {
			rps = int64(MAX_RPS / len(reflect.ValueOf(APIRouteMap).MapKeys()))
		}
		limiter := tollbooth.NewLimiter(rps, time.Second)
		control_methods[ctrl_method](
			fullEndPoint,
			tollbooth_httprouter.LimitHandler(generateAPIEndPoint(handler_method, fullEndPoint), limiter),
		)*/
		control_methods[ctrl_method](
			fullEndPoint,
			generateAPIEndPoint(handler_method, fullEndPoint),
		)
		router.OPTIONS(
			fullEndPoint,
			generateAPIEndPoint(handler_method, fullEndPoint),
		)
	}
	return router
}
