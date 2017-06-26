package api

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"models"
	"net/http"
	"tools"
)

func AddStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	if err := product.AddProductToStoreCategory(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(product)
}

func UpdateStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}

	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	product.UpdateStoreProduct()
	json.NewEncoder(w).Encode(product)
}

func EnableStoreProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product models.Product

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	// set temporary name locally to stop validator from shitting itself
	product.ProductTitle = "temp"
	product.PriceCents = 0
	if validationErr := v.ValidateIncomingJsonRequest(&product); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	product.DB = db.Database
	product.DBSession = product.DB.Session.Copy()
	defer product.DBSession.Close()
	if err := product.ActivateStoreProduct(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(product)
}

func ReorderStoreProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var order models.ProductOrder

	v := new(tools.DefaultValidator)
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		models.WriteError(w, models.ErrBadRequest)
		return
	}
	if validationErr := v.ValidateIncomingJsonRequest(&order); validationErr.Status != 200 {
		models.WriteError(w, &validationErr)
		return
	}
	order.DB = db.Database
	order.DBSession = order.DB.Session.Copy()
	defer order.DBSession.Close()
	if err := order.ReorderStoreProducts(); err != nil {
		models.WriteError(w, models.ErrBadRequest)
	}
	json.NewEncoder(w).Encode(order)
}
