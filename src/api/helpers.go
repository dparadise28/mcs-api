package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"tools"
)

func PaymentMethods(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.PaymentMethods)
}

func OrderMethods(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.OrderMethod)
}

func AllOrderStatuses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.OrderStatuses)
}

func DeliveryOrderStatuses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.DeliveryOrderStatuses)
}

func PickupOrderStatuses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.PickupOrderStatuses)
}

func StatusTree(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.STATUS_TREE)
}

func StatusPaths(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	json.NewEncoder(w).Encode(models.ALLOWED_STATUS_PATH)
}

func EmailQueueLength(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type emailStatus struct {
		Processing int
		Waiting    int
	}
	json.NewEncoder(w).Encode(emailStatus{
		len(tools.EmailQueue),
		len(tools.Emailch),
	})
}
