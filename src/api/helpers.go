package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
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
